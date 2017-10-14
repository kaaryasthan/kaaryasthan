package auth

import (
	jwtmiddleware "github.com/auth0/go-jwt-middleware"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"github.com/kaaryasthan/kaaryasthan/config"
)

var (
	secretKey []byte
)

// JwtMiddleware is middleware to handle all request
var JwtMiddleware = jwtmiddleware.New(jwtmiddleware.Options{
	ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
		//log.Printf("Token: %+v", token)
		return secretKey, nil
	},
	// When set, the middleware verifies that tokens are signed with the specific signing algorithm
	// If the signing method is not constant the ValidationKeyGetter callback can be used to implement additional checks
	// Important to avoid security issues described here: https://auth0.com/blog/2015/03/31/critical-vulnerabilities-in-json-web-token-libraries/
	SigningMethod: jwt.SigningMethodHS256,
})

// Login represents a logged in user session
type Login struct {
	ID       string `jsonapi:"primary,logins"`
	Username string `jsonapi:"attr,username"`
	Password string `jsonapi:"attr,password,omitempty"`
	Token    string `jsonapi:"attr,token,omitempty"`
}

// Register handlers
func Register(art, urt *mux.Router) {
	urt.HandleFunc("/api/v1/register", registerHandler).Methods("POST")
	urt.HandleFunc("/api/v1/login", loginHandler).Methods("POST")
}

func init() {
	secretKey = []byte(config.Config.TokenSecretKey)
}
