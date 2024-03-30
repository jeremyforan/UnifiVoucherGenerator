package UnifiVoucherGenerator

import (
	"fmt"
	"log/slog"
	"net/http"
)

func (c *Client) requestLogin() error {

	urlLogin := c.Url.String() + unifiApiLogin
	urlReferer := c.Url.String() + unifiApiLoginReferer

	req, err := http.NewRequest(http.MethodPost, urlLogin, c.Credentials.HttpPayload())
	if err != nil {
		slog.Error("error creating request login request", "error", err)
		return err
	}

	// Set headers as per the curl command
	addBasicHeaders(req)
	req.Header.Set("Referer", urlReferer)

	body, cookies, err := c.makeRequest(req)

	if !loggedIn(body) {
		err = fmt.Errorf("login failed")
		slog.Error("login failed", "error", err)

		return err
	}

	// todo: move this to another function
	f := false
	for _, cookie := range cookies {
		if cookie.Name == "csrf_token" {
			c.token = cookie.Value
			f = true
			break
		}
	}

	if !f {
		err = fmt.Errorf("csrf_token not found")
		slog.Error("csrf_token not found", "error", err)
		return err
	}

	return nil
}

func (c *Client) requestSelf() (string, error) {
	urlSelf := c.Url.String() + unifiApiSelf
	urlSelfReferer := c.Url.String() + unifiApiLoginReferer

	req, err := http.NewRequest(http.MethodGet, urlSelf, nil)
	if err != nil {
		slog.Error("error creating get self request", "error", err)
		return "", err
	}

	// Set headers as per the curl command
	addBasicHeaders(req)

	req.Header.Set("Referer", urlSelfReferer)
	req.Header.Set("X-Csrf-Token", c.token)

	body, _, err := c.makeRequest(req)
	if err != nil {
		slog.Error("error making request", "error", err)
		return "", err
	}
	return body, nil
}

func (c *Client) requestAddVoucher() error {
	urlVoucher := c.Url.String() + unifiApiCreateVoucher
	urlVoucherReferer := c.Url.String() + unifiApiVoucherReferer

	payload := c.Voucher.HttpPayload()

	req, err := http.NewRequest(http.MethodPost, urlVoucher, payload)
	if err != nil {
		slog.Error("error creating add voucher request", "error", err)
		return err
	}

	addBasicHeaders(req)
	req.Header.Set("Referer", urlVoucherReferer)
	req.Header.Set("X-Csrf-Token", c.token)

	body, _, err := c.makeRequest(req)
	if err != nil {
		slog.Error("error making request", "error", err)
		return err
	}

	// todo: maybe rename this function
	nv, err := processNewVoucherRequestResponse(body)
	if err != nil {
		slog.Error("error processing new voucher request response", "error", err)
		return err
	}

	if !nv.successful() {
		err = fmt.Errorf("voucher request failed")
		slog.Error("voucher request failed", "error", err)
		return err
	}
	return nil
}

func (c *Client) requestFetchPublishedVouchers() (UnifiVouchers, error) {
	urlFetchVouchers := c.Url.String() + unifiApiVouchers
	urlFetchVouchersReferer := c.Url.String() + unifiApiVoucherReferer

	req, err := http.NewRequest(http.MethodPost, urlFetchVouchers, nil)
	if err != nil {
		slog.Error("error creating fetch published vouchers request", "error", err)
		return nil, err
	}

	// Set headers as per the curl command
	addBasicHeaders(req)

	req.Header.Set("Referer", urlFetchVouchersReferer)
	req.Header.Set("X-Csrf-Token", c.token)

	body, _, err := c.makeRequest(req)
	if err != nil {
		slog.Error("error making request", "error", err)
		return nil, err
	}

	vouchers, err := processVoucherListResponse(body)
	if err != nil {
		return nil, err
	}

	return vouchers, nil
}
