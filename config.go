package main

import (
	"gopkg.in/alecthomas/kingpin.v2"
)

// Config project configuration
type Config struct {
	APIVersion string
	TokenFile  string
}

// NewConfig returns new config object storing main configuration options
func NewConfig() *Config {
	return &Config{}
}

// ParseFlags parses passed flags
func (cfg *Config) ParseFlags() {
	cfg.APIVersion = *kingpin.Flag("api-version", "api-version to use").Default("v1beta1").Required().String()
	cfg.TokenFile = *kingpin.Flag("token-file", "csv filepath with prepopulated tokens maped to users in format: `token` `name` `uid`").String()

}
