package UnifiVoucherGenerator

import (
	"io"
	"log/slog"
	"net/http"
)

func addBasicHeaders(req *http.Request) {
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header.Set("Pragma", "no-cache")
	req.Header.Set("DNT", "1") // Do Not Track
}

//todo: may need to be moved
//req.Header.Set("Origin", unifiApiBaseUrl)

func addSecurityHeaders(req *http.Request) {
	req.Header.Add("sec-ch-ua", "\"Chromium\";v=\"122\", \"Not(A:Brand\";v=\"24\", \"Brave\";v=\"122\"")
	req.Header.Add("sec-ch-ua-mobile", "?0")
	req.Header.Add("sec-ch-ua-platform", "\"macOS\"")
	req.Header.Add("sec-fetch-dest", "empty")
	req.Header.Add("sec-fetch-mode", "cors")
	req.Header.Add("sec-fetch-site", "same-origin")
	req.Header.Add("sec-gpc", "1")
}

func loggedIn(responseBody string) bool {
	loginResponse, err := processLoginResponse(responseBody)
	if err != nil {
		return false
	}

	if loginResponse.Meta.Rc == "ok" {
		return true
	}
	return false
}

// makeRequest is a helper function to make a request.go and return the body and cookies
func (c *Client) makeRequest(req *http.Request) (string, []*http.Cookie, error) {
	res, err := c.browser.Do(req)
	if err != nil {
		return "", nil, err
	}

	defer func() {
		err = res.Body.Close()
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

func (v RequestNewVoucherResponse) successful() bool {
	return v.Meta.Rc == "ok"
}
