package main

import (
	_ "figleafproxy/plugin/forwardproxy"
	_ "figleafproxy/plugin/geoip"
	"flag"
	"github.com/mholt/caddy"
	_ "github.com/mholt/caddy/caddyhttp"
	"io/ioutil"
	"log"
	"os"
)

var caddyFile string

func init() {
	flag.StringVar(&caddyFile, "conf", caddy.DefaultConfigFile, "config Caddyfile")
	flag.Parse()

	caddy.SetDefaultCaddyfileLoader("default", caddy.LoaderFunc(defaultLoader))
}

func main() {
	caddy.AppName = "FigLeaf Proxy"
	caddy.AppVersion = "1.0"

	// load caddyfile
	caddyfile, err := caddy.LoadCaddyfile("http")
	if err != nil {
		log.Fatal(err)
	}

	// start caddy server
	instance, err := caddy.Start(caddyfile)
	if err != nil {
		log.Fatal(err)
	}

	instance.Wait()
}

// provide loader function
func defaultLoader(serverType string) (caddy.Input, error) {
	contents, err := ioutil.ReadFile(caddyFile)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	return caddy.CaddyfileInput{
		Contents:       contents,
		Filepath:       caddy.DefaultConfigFile,
		ServerTypeName: serverType,
	}, nil
}
