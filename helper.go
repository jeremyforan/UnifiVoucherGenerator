package UnifiVoucherGenerator

import (
	"io"
	"log/slog"
	"net/http"
	"net/url"
)

// addBasicHeaders adds the basic headers for the http request.
func addBasicHeaders(req *http.Request) {
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header.Set("Pragma", "no-cache")
	req.Header.Set("DNT", "1") // Do Not Track
}

// loggedIn returns true if the login was successful.
func loggedIn(responseBody string) bool {
	loginResponse, err := processLoginResponse(responseBody)
	if err != nil {
		slog.Error("error processing login response", "error", err)
		return false
	}

	if loginResponse.Meta.Rc == "ok" {
		return true
	}

	slog.Error("login failed", "response", responseBody)
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

// loginUrls returns the urls for the login and referer
func (c *Client) loginUrls() (string, string) {
	return c.urlBuilder(unifiApiLogin, unifiApiLoginReferer)
}

// addVoucherUrls returns the urls for the add voucher and referer
func (c *Client) addVoucherUrls() (string, string) {
	return c.urlBuilder(unifiApiCreateVoucher, unifiApiVoucherReferer)
}

// fetchVouchersUrl returns the urls for the fetch vouchers and referer
func (c *Client) fetchVouchersUrl() (string, string) {
	return c.urlBuilder(unifiApiVouchers, unifiApiVoucherReferer)
}

// urlBuilder returns the urls for the endpoint and referer
func (c *Client) urlBuilder(endpoint string, referer string) (string, string) {

	a, err := url.JoinPath(c.Url.String(), endpoint)
	if err != nil {
		slog.Error("error joining url", "error", err)
		return "", ""
	}

	b, err := url.JoinPath(c.Url.String(), referer)
	if err != nil {
		slog.Error("error joining url", "error", err)
		return "", ""
	}

	return a, b
}
