package unifi

import (
	"fmt"
	"io"
	"net/http"
	"strings"
)

type UnifiCredentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Remember bool   `json:"remember"`
	Strict   bool   `json:"strict"`
}

func NewUnifiCredentials(username string, password string) UnifiCredentials {
	return UnifiCredentials{
		Username: username,
		Password: password,
		Remember: true,
		Strict:   true,
	}
}

func (c *Client) Login() error {

	req, err := http.NewRequest(http.MethodPost, unifiApiLogin, c.Credentials.HttpPayload())
	if err != nil {
		return err
	}

	// Set headers as per the curl command
	addBasicHeaders(req)

	req.Header.Set("Referer", unifiApiLoginReferer)

	body, cookies, err := c.MakeRequest(req)
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

func (u UnifiCredentials) String() string {
	return fmt.Sprintf(`{"username":"%s","password":"%s","remember":%t,"strict":%t}`, u.Username, u.Password, u.Remember, u.Strict)
}

func (u UnifiCredentials) HttpPayload() *strings.Reader {
	return strings.NewReader(u.String())
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
