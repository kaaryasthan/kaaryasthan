package config

import (
	"fmt"
	"log"

	"github.com/kelseyhightower/envconfig"
)

// configuration represent all configurations
type configuration struct {
	PostgresHost       string `default:"localhost"`
	PostgresPort       int    `default:"5432"`
	PostgresUser       string `default:"postgres"`
	PostgresPassword   string `default:"secret"`
	PostgresSSLMode    string `default:"disable"`
	KaaryasthanAddress string `default:":8080"`
}

func (c *configuration) PostgresConfig() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s sslmode=%s",
		c.PostgresHost,
		c.PostgresPort,
		c.PostgresUser,
		c.PostgresPassword,
		c.PostgresSSLMode,
	)
}

// Config represent all configurations
var Config configuration

func init() {
	err := envconfig.Process("kaaryasthan", &Config)
	if err != nil {
		log.Fatal(err.Error())
	}
}
