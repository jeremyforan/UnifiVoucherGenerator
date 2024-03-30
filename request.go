package UnifiVoucherGenerator

import (
	"io"
	"log/slog"
	"net/http"
)

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
