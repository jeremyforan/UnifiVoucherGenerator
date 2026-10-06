package UnifiVoucherGenerator

import "fmt"

// FetchVouchers returns every voucher currently known to the controller. Login must have
// been called first.
func (c *Client) FetchVouchers() (UnifiVouchers, error) {
	if c.token == "" {
		return nil, ErrNotLoggedIn
	}
	return c.requestFetchPublishedVouchers()
}

// getVoucherByID finds the voucher whose note matches id. Vouchers created by this library
// carry their Voucher.Id in the note field.
func (v UnifiVouchers) getVoucherByID(id string) (UnifiVoucher, error) {
	for _, vouch := range v {
		if vouch.Note == id {
			return vouch, nil
		}
	}
	return UnifiVoucher{}, fmt.Errorf("%w: no voucher with note %q", ErrVoucherNotFound, id)
}
