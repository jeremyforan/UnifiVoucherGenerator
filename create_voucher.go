package unifi

import (
	"fmt"
	"net/http"
	"strings"
)

type VoucherRequest struct {
	Quota        int    `json:"quota"`
	Note         string `json:"note"`
	Amount       int    `json:"n"`
	ExpireNumber string `json:"expire_number"`
	ExpireUnit   int    `json:"expire_unit"`
	Cmd          string `json:"cmd"`
}

func NewDefaultVoucherRequest() VoucherRequest {
	return VoucherRequest{
		Quota:        5,
		Note:         NoteTimestamp(),
		Amount:       1,
		ExpireNumber: "96",
		ExpireUnit:   60,
		Cmd:          "create-voucher",
	}
}

func (c Client) RequestNewVoucher() error {

	v := NewDefaultVoucherRequest()

	req, err := http.NewRequest(http.MethodPost, unifiApiCreateVoucher, v.HttpPayload())
	if err != nil {
		return err
	}

	addBasicHeaders(req)
	req.Header.Set("Referer", unifiApiVoucherReferer)
	req.Header.Set("DNT", "1") // Do Not Track
	req.Header.Set("X-Csrf-Token", c.token)

	body, _, err := c.MakeRequest(req)
	if err != nil {
		return err
	}
	nv, err := processNewVoucherRequestResponse(body)

	if err != nil {
		return err
	}

	if !nv.successful() {
		return fmt.Errorf("voucher request failed")
	}
	return nil
}

func (v VoucherRequest) String() string {
	return fmt.Sprintf(`{"quota":%d,"note":"%s","n":%d,"expire_number":"%s","expire_unit":%d,"cmd":"%s"}`, v.Quota, v.Note, v.Amount, v.ExpireNumber, v.ExpireUnit, v.Cmd)
}

func (v VoucherRequest) HttpPayload() *strings.Reader {
	return strings.NewReader(v.String())
}

func (v RequestNewVoucherResponse) successful() bool {
	return v.Meta.Rc == "ok"
}
