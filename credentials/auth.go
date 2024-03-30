package credentials

import (
	"fmt"
	"log/slog"
	"strings"
)

// Credentials is a struct that holds the username and password for the Unifi controller
type Credentials struct {
	Username     string `json:"username"`
	Password     string `json:"password"`
	Remember     bool   `json:"remember"`
	Strict       bool   `json:"strict"`
	hidePassword bool
}

// NewCredentials creates a new Credentials struct to be used to log into the Unifi controller
func NewCredentials(username string, password string) Credentials {
	return Credentials{
		Username:     username,
		Password:     password,
		Remember:     true,
		Strict:       true,
		hidePassword: true,
	}
}

// String returns the Credentials struct as a string
func (u Credentials) String() string {
	//todo: this should be a json.Marshal
	return fmt.Sprintf(`{"username":"%s","password":"%s","remember":%t,"strict":%t}`, u.Username, u.Password, u.Remember, u.Strict)
}

// HttpPayload returns the Credentials struct as a strings.Reader to be used in an http request.go as the body
func (u Credentials) HttpPayload() *strings.Reader {
	return strings.NewReader(u.String())
}

// LogValue returns the Credentials struct as a slog.Value for logging.
func (u Credentials) LogValue() slog.Value {

	p := u.Password
	if u.hidePassword {
		p = strings.Repeat("*", len(u.Password))
	}

	return slog.GroupValue(
		slog.String("username", u.Username),
		slog.String("password", p),
		slog.Bool("remember", u.Remember),
		slog.Bool("strict", u.Strict),
	)
}
