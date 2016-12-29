package auth

import (
	"fmt"
	"log"
	"net/http"

	jwtmiddleware "github.com/auth0/go-jwt-middleware"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/kaaryasthan/kaaryasthan/config"
	"github.com/kaaryasthan/kaaryasthan/route"
)

var (
	privateKey []byte
	publicKey  []byte
)

// OAuth2 represents a OAuth 2 provider
type OAuth2 struct {
	Name string
}

// Register a OAuth 2 provider
func Register(name string, begin func(http.ResponseWriter, *http.Request), complete func(http.ResponseWriter, *http.Request)) *OAuth2 {
	o := OAuth2{Name: name}
	route.URT.HandleFunc(fmt.Sprintf("/api/v1/auth/%s", name), begin).Methods("GET")
	route.URT.HandleFunc(fmt.Sprintf("/api/v1/auth/%s/callback", name), complete).Methods("GET")
	return &o
}

// JwtMiddleware is middleware to handle all request
var JwtMiddleware = jwtmiddleware.New(jwtmiddleware.Options{
	ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
		log.Printf("Token: %+v", token)
		return publicKey, nil
	},
	// When set, the middleware verifies that tokens are signed with the specific signing algorithm
	// If the signing method is not constant the ValidationKeyGetter callback can be used to implement additional checks
	// Important to avoid security issues described here: https://auth0.com/blog/2015/03/31/critical-vulnerabilities-in-json-web-token-libraries/
	SigningMethod: jwt.SigningMethodRS256,
})

func init() {
	// FIXME: Verify key
	privateKey = []byte(config.Config.TokenPrivateKey)
	publicKey = []byte(config.Config.TokenPublicKey)
}
