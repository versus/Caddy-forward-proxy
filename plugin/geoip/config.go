package geoip

import (
	"errors"
	"github.com/mholt/caddy"
)

// Config specifies configuration parsed for Caddyfile
type Config struct {
	DatabasePath string
}

func parseConfig(c *caddy.Controller) (*Config, error) {
	for c.Next() {
		if c.Val() != "geoip" {
			continue
		}

		if !c.NextArg() {
			continue
		}

		return &Config{c.Val()}, nil
	}
	return nil, errors.New("geoip directive not found")
}
