package unifi

import "encoding/json"

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
