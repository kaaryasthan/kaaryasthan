package auth

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	jwtmiddleware "github.com/auth0/go-jwt-middleware"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/kaaryasthan/kaaryasthan/route"
	"github.com/urfave/negroni"
)

var (
	privateKey []byte
	publicKey  []byte
)

// Facebook OAuth 2
//_ "github.com/kaaryasthan/kaaryasthan/auth/facebook"
// Google OAuth 2
//_ "github.com/kaaryasthan/kaaryasthan/auth/google"

// OAuth2 represents a OAuth 2 provider
type OAuth2 struct {
	Name string
}

// Register a OAuth 2 provider
func Register(name string) *OAuth2 {
	o := OAuth2{Name: name}
	return &o
}

// Token represents a token
type Token struct {
	Success bool   `json:"success"`
	Token   string `json:"token"`
	Message string `json:"message"`
}

func authHandler(w http.ResponseWriter, r *http.Request) {
	token := jwt.New(jwt.GetSigningMethod("RS256"))
	token.Claims.(jwt.MapClaims)["sub"] = "guest"
	token.Claims.(jwt.MapClaims)["exp"] = time.Now().Add(time.Hour * 24 * 7).Unix()
	tokenString, _ := token.SignedString(privateKey)
	log.Printf("Valid Token: %+v", token)
	log.Printf("tokenString: %v\n", tokenString)

	authToken, err := json.Marshal(Token{true, tokenString, "Logged in"})
	if err != nil {
		log.Fatal("Unable to marhal token")
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(authToken))
}

func init() {
	jwtMiddleware := jwtmiddleware.New(jwtmiddleware.Options{
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			log.Printf("Token: %+v", token)
			return publicKey, nil
		},
	})

	route.URT.HandleFunc("/api/v1/auth", authHandler).Methods("POST")
	route.URT.PathPrefix("/api").Handler(
		negroni.New(negroni.HandlerFunc(jwtMiddleware.HandlerWithNext), negroni.Wrap(route.RT)))
}
