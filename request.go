package UnifiVoucherGenerator

import (
	"fmt"
	"net/http"
)

func (c *Client) requestLogin() error {
	urlLogin, urlReferer := c.loginUrls()

	req, err := http.NewRequest(http.MethodPost, urlLogin, c.Credentials.HttpPayload())
	if err != nil {
		return fmt.Errorf("creating login request: %w", err)
	}

	addBasicHeaders(req)
	req.Header.Set("Referer", urlReferer)

	body, cookies, err := c.makeRequest(req)
	if err != nil {
		return fmt.Errorf("login request: %w", err)
	}

	if !loggedIn(body) {
		return fmt.Errorf("%w: %s", ErrLoginFailed, responseMessage(body))
	}

	token, ok := csrfTokenFromCookies(cookies)
	if !ok {
		return ErrCSRFTokenNotFound
	}
	c.token = token

	return nil
}

func (c *Client) requestAddVoucher() error {
	urlVoucher, urlVoucherReferer := c.addVoucherUrls()

	payload := c.Voucher.HttpPayload()
	if payload == nil {
		return fmt.Errorf("encoding voucher payload failed")
	}

	req, err := http.NewRequest(http.MethodPost, urlVoucher, payload)
	if err != nil {
		return fmt.Errorf("creating add voucher request: %w", err)
	}

	addBasicHeaders(req)
	req.Header.Set("Referer", urlVoucherReferer)
	req.Header.Set("X-Csrf-Token", c.token)

	body, _, err := c.makeRequest(req)
	if err != nil {
		return fmt.Errorf("add voucher request: %w", err)
	}

	nv, err := processNewVoucherRequestResponse(body)
	if err != nil {
		return fmt.Errorf("decoding add voucher response: %w", err)
	}

	if !nv.successful() {
		return fmt.Errorf("%w: %s", ErrVoucherRequestFailed, responseMessage(body))
	}
	return nil
}

func (c *Client) requestFetchPublishedVouchers() (UnifiVouchers, error) {
	urlFetchVouchers, urlFetchVouchersReferer := c.fetchVouchersUrl()

	req, err := http.NewRequest(http.MethodPost, urlFetchVouchers, nil)
	if err != nil {
		return nil, fmt.Errorf("creating fetch vouchers request: %w", err)
	}

	addBasicHeaders(req)
	req.Header.Set("Referer", urlFetchVouchersReferer)
	req.Header.Set("X-Csrf-Token", c.token)

	body, _, err := c.makeRequest(req)
	if err != nil {
		return nil, fmt.Errorf("fetch vouchers request: %w", err)
	}

	vouchers, err := processVoucherListResponse(body)
	if err != nil {
		return nil, fmt.Errorf("decoding fetch vouchers response: %w", err)
	}

	return vouchers, nil
}
