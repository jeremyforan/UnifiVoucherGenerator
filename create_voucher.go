package UnifiVoucherGenerator

import (
	"fmt"
	"github.com/jeremyforan/UnifiVoucherGenerator/voucher"
	"net/http"
	"strings"
)

func (c Client) RequestNewVoucher() error {

	v := voucher.NewDefaultVoucherRequest()

	req, err := http.NewRequest(http.MethodPost, unifiApiCreateVoucher, v.HttpPayload())
	if err != nil {
		return err
	}

	addBasicHeaders(req)
	req.Header.Set("Referer", unifiApiVoucherReferer)
	req.Header.Set("DNT", "1") // Do Not Track
	req.Header.Set("X-Csrf-Token", c.token)

	body, _, err := c.buildRequest(req)
	if err != nil {
		return err
	}
	nv, err := processNewVoucherRequestResponse(body)

	if err != nil {
		return err
	}

	if !nv.successful() {
		return fmt.Errorf("voucher request.go failed")
	}
	return nil
}

func (v voucher.VoucherRequest) String() string {
	return fmt.Sprintf(`{"quota":%d,"note":"%s","n":%d,"expire_number":"%s","expire_unit":%d,"cmd":"%s"}`, v.Quota, v.Note, v.Amount, v.ExpireNumber, v.ExpireUnit, v.Cmd)
}

func (v voucher.VoucherRequest) HttpPayload() *strings.Reader {
	return strings.NewReader(v.String())
}

func (v RequestNewVoucherResponse) successful() bool {
	return v.Meta.Rc == "ok"
}
