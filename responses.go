package UnifiVoucherGenerator

import "encoding/json"

//todo: add logging

type Meta struct {
	Rc string `json:"rc"` // Maps the "rc" field to check if it's "ok"
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

func processResponse[T any](body string) (*T, error) {
	var response T

	err := json.Unmarshal([]byte(body), &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

// ProcessLoginResponse converts the JSON response from the login endpoint into a struct
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

func processVoucherListResponse(body string) (UnifiVouchers, error) {
	t, err := processResponse[VoucherListResponse](body)
	if err != nil {
		return UnifiVouchers{}, err
	}
	return t.Data, nil
}

// VoucherListResponse Define struct for the top-level JSON object
type VoucherListResponse struct {
	Meta Meta          `json:"meta"`
	Data UnifiVouchers `json:"data"`
}
