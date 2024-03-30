package UnifiVoucherGenerator

import (
	"errors"
	"log/slog"
	"net/http"
)

func (c *Client) FetchVouchers() (UnifiVouchers, error) {
	urlFetchVouchers := c.Url.String() + unifiApiVouchers
	urlFetchVouchersReferer := c.Url.String() + unifiApiVoucherReferer

	req, err := http.NewRequest(http.MethodPost, urlFetchVouchers, nil)
	if err != nil {
		return []UnifiVoucher{}, err
	}

	// Set headers as per the curl command
	addBasicHeaders(req)

	req.Header.Set("Referer", urlFetchVouchersReferer)
	req.Header.Set("X-Csrf-Token", c.token)

	body, _, err := c.makeRequest(req)

	vouchers, err := processVoucherListResponse(body)
	if err != nil {
		return []UnifiVoucher{}, err
	}

	return vouchers, nil
}

func (v UnifiVouchers) getVoucherByID(id string) (UnifiVoucher, error) {

	for _, vouch := range v {
		if vouch.Note == id {
			return vouch, nil
		}
	}

	err := errors.New("voucher not found")
	slog.Error("voucher not found", "error", err)
	return UnifiVoucher{}, err
}
