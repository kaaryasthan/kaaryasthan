package google

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/kaaryasthan/kaaryasthan/auth"
	"github.com/kaaryasthan/kaaryasthan/config"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var (
	privateKey []byte
	publicKey  []byte
)

var googleconf = &oauth2.Config{
	ClientID:     config.Config.GoogleClientID,
	ClientSecret: config.Config.GoogleClientSecret,
	RedirectURL:  config.Config.GoogleRedirectURL,
	Scopes: []string{
		"openid",
		"profile",
		"email",
	},
	Endpoint: google.Endpoint,
}

func completeAuthHandler(w http.ResponseWriter, r *http.Request) {
	authcode := r.FormValue("code")

	tok, err := googleconf.Exchange(oauth2.NoContext, authcode)
	if err != nil {
		fmt.Println("err is", err)
	}

	fmt.Println("token is ", tok)
	response, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + tok.AccessToken)

	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)
	log.Printf(string(contents))
}

// Token represents a token
type Token struct {
	Success bool   `json:"success"`
	Token   string `json:"token"`
	Message string `json:"message"`
}

func beginAuthHandler(w http.ResponseWriter, r *http.Request) {
	url := googleconf.AuthCodeURL("state")
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
	return

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
	// FIXME: Verify key
	privateKey = []byte(config.Config.TokenPrivateKey)
	publicKey = []byte(config.Config.TokenPublicKey)
	auth.Register("google", beginAuthHandler, completeAuthHandler)
}
