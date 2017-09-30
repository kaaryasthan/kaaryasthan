package auth

import (
	"crypto/rand"
	"io"
	"log"

	jwtmiddleware "github.com/auth0/go-jwt-middleware"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"github.com/kaaryasthan/kaaryasthan/config"
	"github.com/kaaryasthan/kaaryasthan/jsonapi"
)

var (
	privateKey []byte
	publicKey  []byte
)

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

// Schema represents a database schema
type Schema struct {
	Username    string
	Name        string
	Email       string
	Password    string
	AccessToken *string
}

// New returns a schema
func New(d jsonapi.Data) *Schema {
	s := &Schema{}
	s.Username = d.Attributes["username"]
	s.Name = d.Attributes["name"]
	s.Email = d.Attributes["email"]
	s.Password = d.Attributes["password"]
	return s
}

func randomSalt() []byte {
	s := make([]byte, 32)
	_, err := io.ReadFull(rand.Reader, s)
	if err != nil {
		log.Fatal(err)
	}
	return s
}

// Register handlers
func Register(art, urt *mux.Router) {
	urt.HandleFunc("/api/v1/register", registerHandler).Methods("POST")
	urt.HandleFunc("/api/v1/login", loginHandler).Methods("POST")
}

func init() {
	// FIXME: Verify key
	privateKey = []byte(config.Config.TokenPrivateKey)
	publicKey = []byte(config.Config.TokenPublicKey)
}
