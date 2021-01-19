package config

import (
	"flag"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// Config is app of conig
type Config struct {
	// Network parameters. Set mainnet, testnet, or regtest using this.
	RESTConfig RESTConfig `mapstructure:"rest" json:"rest"`
	WSConfig   WSConfig   `mapstructure:"ws" json:"ws"`
}

type RESTConfig struct {
	ConnAddr   string `mapstructure:"connect" json:"connect"`
	ListenAddr string `mapstructure:"listen" json:"listen"`
}

type WSConfig struct {
	ListenAddr string `mapstructure:"listen" json:"listen"`
}

func init() {
	// Bind rest flags
	pflag.StringP("rest.listen", "l", "0.0.0.0:9090", "The listen address for REST API")
	// Bind ws flags
	pflag.StringP("ws.listen", "w", "0.0.0.0:9091", "The listen address for Websocket API")
}

// NewDefaultConfig is default config
func NewDefaultConfig() (*Config, error) {
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	pflag.Parse()
	viper.BindPFlags(pflag.CommandLine)
	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}
	return &config, nil
}
