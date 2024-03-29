package unifi

import (
	"net/http"
)

// VoucherListResponse Define struct for the top-level JSON object
type VoucherListResponse struct {
	Meta Meta          `json:"meta"`
	Data UnifiVouchers `json:"data"`
}

// UnifiVoucher Define struct for each item in the data array
type UnifiVoucher struct {
	Duration      int    `json:"duration"`
	QosOverwrite  bool   `json:"qos_overwrite"`
	Note          string `json:"note"`
	Code          string `json:"code"`
	ForHotspot    bool   `json:"for_hotspot"`
	CreateTime    int64  `json:"create_time"`
	Quota         int    `json:"quota"`
	SiteID        string `json:"site_id"`
	ID            string `json:"_id"`
	AdminName     string `json:"admin_name"`
	Used          int    `json:"used"`
	Status        string `json:"status"`
	StatusExpires int    `json:"status_expires"`
}

type UnifiVouchers []UnifiVoucher

func (c *Client) FetchVouchers() (UnifiVouchers, error) {

	req, err := http.NewRequest(http.MethodPost, unifiApiVouchers, nil)
	if err != nil {
		return []UnifiVoucher{}, err
	}

	// Set headers as per the curl command
	addBasicHeaders(req)

	req.Header.Set("Referer", unifiApiVoucherReferer)
	req.Header.Set("X-Csrf-Token", c.token)

	body, _, err := c.MakeRequest(req)

	vouchers, err := processVoucherListResponse(body)
	if err != nil {
		return []UnifiVoucher{}, err
	}

	return vouchers, nil
}

func (v UnifiVouchers) GetLatestVoucher() UnifiVoucher {
	latestTS := int64(0)
	newestVoucher := UnifiVoucher{}

	for _, voucher := range v {
		if voucher.Status == "VALID_MULTI" {
			if voucher.CreateTime > latestTS {
				latestTS = voucher.CreateTime
				newestVoucher = voucher
			}
		}
	}

	return newestVoucher
}

func (v UnifiVoucher) GetVoucherByCode() (VoucherCode, error) {
	c, err := NewVoucherFromString(v.Code)
	if err != nil {
		return VoucherCode{}, err
	}
	return c, nil

}
