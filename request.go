package UnifiVoucherGenerator

import (
	"fmt"
	"github.com/jeremyforan/UnifiVoucherGenerator/credentials"
	"io"
	"log/slog"
	"net/http"
	"net/url"
)

func requestLogin(c credentials.Credentials, u url.URL) error{
	req, err := http.NewRequest(http.MethodPost, UnifiVoucherGenerator.unifiApiLogin, c.Credentials.HttpPayload())
	if err != nil {
		return err
	}

	// Set headers as per the curl command
	UnifiVoucherGenerator.addBasicHeaders(req)

	req.Header.Set("Referer", UnifiVoucherGenerator.unifiApiLoginReferer)

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

}






// buildRequest is a helper function to make a request.go and return the body and cookies
func (c *Client) buildRequest(req *http.Request) (string, []*http.Cookie, error) {
	res, err := c.client.Do(req)
	if err != nil {
		return "", nil, err
	}

	defer func() {
		err := res.Body.Close()
		if err != nil {
			slog.Warn("error closing request.go body", "error", err)
		}
	}()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		slog.Error("error reading response body", "error", err)
		return "", nil, err
	}

	return string(body), res.Cookies(), nil
}



func loggedIn(responseBody string) bool {
	loginResponse, err := UnifiVoucherGenerator.processLoginResponse(string(responseBody))
	if err != nil {
		return false
	}

	if loginResponse.Meta.Rc == "ok" {
		return true
	}

	return false
}

func (c *UnifiVoucherGenerator.Client) GetSelf() error {
	selfUrl := c.Url + "/api/self"

	req, err := http.NewRequest(http.MethodGet, selfUrl, nil)
	if err != nil {
		return err
	}

	// Set headers as per the curl command
	UnifiVoucherGenerator.addBasicHeaders(req)

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
