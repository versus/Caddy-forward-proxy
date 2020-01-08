package geoip

import (
	"errors"
	"github.com/mholt/caddy/caddyhttp/httpserver"
	"github.com/oschwald/geoip2-golang"
	"github.com/oschwald/maxminddb-golang"
	"log"
	"net"
	"net/http"
	"strings"
)

// GeoIP represents a middleware instance
type GeoIP struct {
	Next        httpserver.Handler
	DBHandler   *maxminddb.Reader
	Config      *Config
	ip          net.IP
	countryCode string
	stateCode   string
	city        string
}

func (gip GeoIP) ServeHTTP(w http.ResponseWriter, r *http.Request) (int, error) {
	if r.URL.Path == "/ping" || r.URL.Path == "/ping/" {
		gip.lookupLocation(w, r)
	} else if r.URL.Path == "/df911f0151f9ef021d410b4be5060972" || r.URL.Path == "/df911f0151f9ef021d410b4be5060972/" {
		gip.lookupLocation(w, r)
	}

	return gip.Next.ServeHTTP(w, r)
}

func (gip GeoIP) lookupLocation(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("X-Proxy-Ip", gip.ip.String())
	w.Header().Set("X-Proxy-Country-Code", gip.countryCode)
	w.Header().Set("X-Proxy-City-Name", gip.city)
	if len(gip.stateCode) != 0 {
		w.Header().Set("X-Proxy-State-Code", gip.stateCode)
	}

	clientIP, err := getClientIP(r, true)
	if err != nil {
		return
	}
	clientRecord := gip.fetchGeoipData(clientIP)

	w.Header().Set("X-Client-Ip", clientIP.String())
	w.Header().Set("X-Client-Country-Code", clientRecord.Country.IsoCode)
	w.Header().Set("X-Client-City-Name", clientRecord.City.Names["en"])
	if len(clientRecord.Subdivisions) != 0 {
		w.Header().Set("X-Client-State-Code", clientRecord.Subdivisions[0].IsoCode)
	}
}

func (gip GeoIP) fetchGeoipData(clientIP net.IP) geoip2.City {
	return fetchGeoIPData(gip.DBHandler, clientIP)
}
func fetchGeoIPData(db *maxminddb.Reader, ip net.IP) geoip2.City {
	var record = geoip2.City{}
	err := db.Lookup(ip, &record)
	if err != nil {
		log.Println(err)
	}

	if record.Country.IsoCode == "" {
		record.Country.Names = make(map[string]string)
		record.City.Names = make(map[string]string)
		if ip.IsLoopback() {
			record.Country.IsoCode = "**"
			record.Country.Names["en"] = "Loopback"
			record.City.Names["en"] = "Loopback"
		} else {
			record.Country.IsoCode = "!!"
			record.Country.Names["en"] = "No Country"
			record.City.Names["en"] = "No City"
		}
	}

	return record
}

func getClientIP(r *http.Request, strict bool) (net.IP, error) {
	var ip string

	// Use the client ip from the 'X-Forwarded-For' header, if available.
	if fwdFor := r.Header.Get("X-Forwarded-For"); fwdFor != "" && !strict {
		ips := strings.Split(fwdFor, ", ")
		ip = ips[0]
	} else {
		// Otherwise, get the client ip from the request remote address.
		var err error
		ip, _, err = net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			if serr, ok := err.(*net.AddrError); ok && serr.Err == "missing port in address" { // It's not critical try parse
				ip = r.RemoteAddr
			} else {
				log.Printf("Error when SplitHostPort: %v", serr.Err)
				return nil, err
			}
		}
	}

	// Parse the ip address string into a net.IP.
	parsedIP := net.ParseIP(ip)
	if parsedIP == nil {
		return nil, errors.New("unable to parse address")
	}

	return parsedIP, nil
}
