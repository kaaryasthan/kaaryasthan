package config

import (
	"fmt"
	"log"

	"github.com/kelseyhightower/envconfig"
)

// Configuration represent all configurations
type Configuration struct {
	PostgresHost     string `default:"localhost" split_words:"true"`
	PostgresPort     int    `default:"5430" split_words:"true"`
	PostgresUser     string `default:"postgres" split_words:"true"`
	PostgresDatabase string `default:"postgres" split_words:"true"`
	PostgresPassword string `default:"secret" split_words:"true"`
	PostgresSSLMode  string `default:"disable" envconfig:"POSTGRES_SSL_MODE"`
	HTTPAddress      string `default:":8080" envconfig:"HTTP_ADDRESS"`
	TokenSecretKey   string `default:"secret" split_words:"true"`
	DeveloperMode    bool   `default:"false" split_words:"true"`
	BleveIndexPath   string `split_words:"true"`
	SendinblueKey    string `default:"secret" split_words:"true"`
	BaseURL          string `envconfig:"BASE_URL"`
	EmailSender      string `default:"noreply@example.org" split_words:"true"`
	EmailReplyTo     string `default:"noreply@example.org" split_words:"true"`
}

// PostgresConfig provides PostgreSQL connection string
func (c *Configuration) PostgresConfig() string {
	return fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s sslmode=%s",
		c.PostgresHost,
		c.PostgresPort,
		c.PostgresUser,
		c.PostgresDatabase,
		c.PostgresPassword,
		c.PostgresSSLMode,
	)
}

// SetDatabaseName set database
func (c *Configuration) SetDatabaseName(dbname string) {
	c.PostgresDatabase = dbname
}

// SetBleveIndexPath set Bleve index path
func (c *Configuration) SetBleveIndexPath(path string) {
	c.BleveIndexPath = path
}

// Config represent all configurations
var Config Configuration

func init() {
	err := envconfig.Process("kaaryasthan", &Config)
	if err != nil {
		log.Fatal(err.Error())
	}
}
