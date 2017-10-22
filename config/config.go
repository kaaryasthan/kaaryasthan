package config

import (
	"fmt"
	"log"

	"github.com/kelseyhightower/envconfig"
)

// configuration represent all configurations
type configuration struct {
	PostgresHost     string `default:"localhost" split_words:"true"`
	PostgresPort     int    `default:"5430" split_words:"true"`
	PostgresUser     string `default:"postgres" split_words:"true"`
	PostgresDatabase string `default:"postgres" split_words:"true"`
	PostgresPassword string `default:"secret" split_words:"true"`
	PostgresSSLMode  string `default:"disable" envconfig:"POSTGRES_SSL_MODE"`
	HTTPAddress      string `default:":8080" envconfig:"HTTP_ADDRESS"`
	TokenSecretKey   string `default:"secret" split_words:"true"`
	DeveloperMode    bool   `default:"false" split_words:"true"`
}

func (c *configuration) PostgresConfig() string {
	return fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s sslmode=%s",
		c.PostgresHost,
		c.PostgresPort,
		c.PostgresUser,
		c.PostgresDatabase,
		c.PostgresPassword,
		c.PostgresSSLMode,
	)
}

func (c *configuration) SetDatabaseName(dbname string) {
	c.PostgresDatabase = dbname
}

// Config represent all configurations
var Config configuration

func init() {
	err := envconfig.Process("kaaryasthan", &Config)
	if err != nil {
		log.Fatal(err.Error())
	}
}
