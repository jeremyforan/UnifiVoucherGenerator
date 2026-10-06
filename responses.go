package UnifiVoucherGenerator

import (
	"encoding/json"
	"fmt"
)

// Meta is the status envelope returned by every controller endpoint.
type Meta struct {
	Rc  string `json:"rc"`            // "ok" on success, "error" otherwise
	Msg string `json:"msg,omitempty"` // error code such as "api.err.LoginRequired", set when Rc is "error"
}

func (m Meta) ok() bool {
	return m.Rc == "ok"
}

type LoginResponse struct {
	Meta Meta          `json:"meta"` // Maps the "meta" field
	Data []interface{} `json:"data"` // Use []interface{} for arbitrary data; adjust as needed
}

type RequestNewVoucherResponse struct {
	Meta struct {
		Rc string `json:"rc"`
	} `json:"meta"`
	Data []struct {
		CreateTime int `json:"create_time"`
	} `json:"data"`
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

// VoucherListResponse Define struct for the top-level JSON object
type VoucherListResponse struct {
	Meta Meta          `json:"meta"`
	Data UnifiVouchers `json:"data"`
}

func processResponse[T any](body string) (*T, error) {
	var response T

	err := json.Unmarshal([]byte(body), &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

// processLoginResponse converts the JSON response from the login endpoint into a struct
func processLoginResponse(body string) (*LoginResponse, error) {
	return processResponse[LoginResponse](body)
}

func processNewVoucherRequestResponse(body string) (RequestNewVoucherResponse, error) {
	t, err := processResponse[RequestNewVoucherResponse](body)
	if err != nil {
		return RequestNewVoucherResponse{}, err
	}
	return *t, nil
}

// processVoucherListResponse decodes the voucher list. A response whose meta.rc is not
// "ok" (for example an expired session) is returned as an error rather than an empty list.
func processVoucherListResponse(body string) (UnifiVouchers, error) {
	t, err := processResponse[VoucherListResponse](body)
	if err != nil {
		return UnifiVouchers{}, err
	}
	if !t.Meta.ok() {
		return UnifiVouchers{}, fmt.Errorf("controller returned %q: %s", t.Meta.Rc, t.Meta.Msg)
	}
	return t.Data, nil
}

// responseMessage extracts a human readable status from a controller response body for
// use in error messages. It falls back to the raw body when the body is not a known shape.
func responseMessage(body string) string {
	var envelope struct {
		Meta Meta `json:"meta"`
	}
	if err := json.Unmarshal([]byte(body), &envelope); err != nil || envelope.Meta.Rc == "" {
		return body
	}
	if envelope.Meta.Msg == "" {
		return envelope.Meta.Rc
	}
	return envelope.Meta.Msg
}
