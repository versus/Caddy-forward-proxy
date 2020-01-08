package geoip

import (
	"github.com/aws/aws-sdk-go-v2/aws/ec2metadata"
	"github.com/aws/aws-sdk-go-v2/aws/external"
	"log"
	"net"
	"time"

	"github.com/mholt/caddy"
	"github.com/mholt/caddy/caddyhttp/httpserver"
	"github.com/oschwald/maxminddb-golang"
)

// Init initializes the plugin
func init() {
	caddy.RegisterPlugin("geoip", caddy.Plugin{
		ServerType: "http",
		Action:     setup,
	})
}

func setup(c *caddy.Controller) error {
	config, err := parseConfig(c)
	if err != nil {
		return err
	}

	dbhandler, err := maxminddb.Open(config.DatabasePath)
	if err != nil {
		return c.Err("geoip: Can't open database: " + config.DatabasePath)
	}

	awsCfg, err := external.LoadDefaultAWSConfig()
	if err != nil {
		log.Fatal(err)
	}

	awsCfg.HTTPClient.Timeout = time.Second * 5
	m := ec2metadata.New(awsCfg)

	var ip string
	ip, err = m.GetMetadata("public-ipv4")
	if err != nil {
		ip = "127.0.0.1"
	}
	serverIP := net.ParseIP(ip)
	serverRecord := fetchGeoIPData(dbhandler, serverIP)
	var stateCode string
	if len(serverRecord.Subdivisions) != 0 {
		stateCode = serverRecord.Subdivisions[0].IsoCode
	}

	// Create new middleware
	newMiddleWare := func(next httpserver.Handler) httpserver.Handler {
		return &GeoIP{
			Next:        next,
			DBHandler:   dbhandler,
			Config:      config,
			ip:          serverIP,
			countryCode: serverRecord.Country.IsoCode,
			stateCode:   stateCode,
			city:        serverRecord.City.Names["en"],
		}
	}
	// Add middleware
	cfg := httpserver.GetConfig(c)
	cfg.AddMiddleware(newMiddleWare)

	return nil
}
