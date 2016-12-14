package auth

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

func init() {
}
