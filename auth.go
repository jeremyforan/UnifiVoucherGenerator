package UnifiVoucherGenerator

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"
)

// UnifiCredentials is a struct that holds the username and password for the Unifi controller
type UnifiCredentials struct {
	Username     string `json:"username"`
	Password     string `json:"password"`
	Remember     bool   `json:"remember"`
	Strict       bool   `json:"strict"`
	hidePassword bool
}

// NewUnifiCredentials creates a new UnifiCredentials struct to be used to log into the Unifi controller
func NewUnifiCredentials(username string, password string) UnifiCredentials {
	return UnifiCredentials{
		Username:     username,
		Password:     password,
		Remember:     true,
		Strict:       true,
		hidePassword: true,
	}
}

// String returns the UnifiCredentials struct as a string
func (u UnifiCredentials) String() string {
	return fmt.Sprintf(`{"username":"%s","password":"%s","remember":%t,"strict":%t}`, u.Username, u.Password, u.Remember, u.Strict)
}

// HttpPayload returns the UnifiCredentials struct as a strings.Reader to be used in an http request.go as the body
func (u UnifiCredentials) HttpPayload() *strings.Reader {
	return strings.NewReader(u.String())
}

// LogValue returns the UnifiCredentials struct as a slog.Value for logging.
func (u UnifiCredentials) LogValue() slog.Value {

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

func loggedIn(responseBody string) bool {
	loginResponse, err := processLoginResponse(string(responseBody))
	if err != nil {
		return false
	}

	if loginResponse.Meta.Rc == "ok" {
		return true
	}

	return false
}

func (c *Client) Login() error {

	req, err := http.NewRequest(http.MethodPost, unifiApiLogin, c.Credentials.HttpPayload())
	if err != nil {
		return err
	}

	// Set headers as per the curl command
	addBasicHeaders(req)

	req.Header.Set("Referer", unifiApiLoginReferer)

	body, cookies, err := c.buildRequest(req)
	for _, cookie := range cookies {
		if cookie.Name == "csrf_token" {
			c.token = cookie.Value
			break
		}
	}

	if loggedIn(string(body)) {
		return nil
	}

	return fmt.Errorf("login failed")
}

func (c *Client) GetSelf() error {
	selfUrl := c.Url + "/api/self"

	req, err := http.NewRequest(http.MethodGet, selfUrl, nil)
	if err != nil {
		return err
	}

	// Set headers as per the curl command
	addBasicHeaders(req)

	req.Header.Set("Referer", c.Url+"/manage/account/login")
	req.Header.Set("DNT", "1") // Do Not Track
	req.Header.Set("X-Csrf-Token", c.token)

	res, err := c.Client.Do(req) // Use the client with a cookie jar
	if err != nil {
		return err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	fmt.Println(string(body))
	return nil
}
