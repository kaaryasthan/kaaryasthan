package config

import (
	"fmt"
	"log"

	"github.com/kelseyhightower/envconfig"
)

// configuration represent all configurations
type configuration struct {
	PostgresHost      string   `default:"localhost" split_words:"true"`
	PostgresPort      int      `default:"5433" split_words:"true"`
	PostgresUser      string   `default:"postgres" split_words:"true"`
	PostgresPassword  string   `default:"secret" split_words:"true"`
	PostgresSSLMode   string   `default:"disable" envconfig:"POSTGRES_SSL_MODE"`
	HTTPAddress       string   `default:":8080" envconfig:"HTTP_ADDRESS"`
	TokenPrivateKey   string   `split_words:"true"`
	TokenPublicKey    string   `split_words:"true"`
	ClientID          string   `envconfig:"CLIENT_ID"`
	ClientSecret      string   `split_words:"true"`
	IdentityProviders []string `default:"google" split_words:"true"`
	RedirectURL       string   `envconfig:"REDIRECT_URL"`
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
