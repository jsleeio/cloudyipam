package shared

import (
	"github.com/spf13/viper"

	"github.com/jsleeio/cloudyipam/pkg/cloudyipam"
)

func DatabaseConfigFromViper() cloudyipam.DatabaseConfig {
	return cloudyipam.DatabaseConfig{
		Database: viper.GetString("db-name"),
		User:     viper.GetString("db-user"),
		Host:     viper.GetString("db-host"),
		Port:     viper.GetString("db-port"),
		Password: viper.GetString("db-password"),
		TLS:      viper.GetBool("db-tls"),
	}
}
