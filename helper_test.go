package UnifiVoucherGenerator

import (
	"github.com/jeremyforan/UnifiVoucherGenerator/credentials"
	"github.com/jeremyforan/UnifiVoucherGenerator/voucher"
	"net/http"
	"net/url"
	"reflect"
	"testing"
)

func TestClient_addVoucherUrls(t *testing.T) {
	type fields struct {
		Credentials credentials.Credentials
		browser     *http.Client
		Url         *url.URL
		token       string
		Voucher     *voucher.Voucher
	}
	var tests []struct {
		name   string
		fields fields
		want   string
		want1  string
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				Credentials: tt.fields.Credentials,
				browser:     tt.fields.browser,
				Url:         tt.fields.Url,
				token:       tt.fields.token,
				Voucher:     tt.fields.Voucher,
			}
			got, got1 := c.addVoucherUrls()
			if got != tt.want {
				t.Errorf("addVoucherUrls() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("addVoucherUrls() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestClient_fetchVouchersUrl(t *testing.T) {
	type fields struct {
		Credentials credentials.Credentials
		browser     *http.Client
		Url         *url.URL
		token       string
		Voucher     *voucher.Voucher
	}
	var tests []struct {
		name   string
		fields fields
		want   string
		want1  string
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				Credentials: tt.fields.Credentials,
				browser:     tt.fields.browser,
				Url:         tt.fields.Url,
				token:       tt.fields.token,
				Voucher:     tt.fields.Voucher,
			}
			got, got1 := c.fetchVouchersUrl()
			if got != tt.want {
				t.Errorf("fetchVouchersUrl() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("fetchVouchersUrl() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestClient_loginUrls(t *testing.T) {
	type fields struct {
		Credentials credentials.Credentials
		browser     *http.Client
		Url         *url.URL
		token       string
		Voucher     *voucher.Voucher
	}
	var tests []struct {
		name   string
		fields fields
		want   string
		want1  string
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				Credentials: tt.fields.Credentials,
				browser:     tt.fields.browser,
				Url:         tt.fields.Url,
				token:       tt.fields.token,
				Voucher:     tt.fields.Voucher,
			}
			got, got1 := c.loginUrls()
			if got != tt.want {
				t.Errorf("loginUrls() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("loginUrls() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestClient_makeRequest(t *testing.T) {
	type fields struct {
		Credentials credentials.Credentials
		browser     *http.Client
		Url         *url.URL
		token       string
		Voucher     *voucher.Voucher
	}
	type args struct {
		req *http.Request
	}
	var tests []struct {
		name    string
		fields  fields
		args    args
		want    string
		want1   []*http.Cookie
		wantErr bool
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				Credentials: tt.fields.Credentials,
				browser:     tt.fields.browser,
				Url:         tt.fields.Url,
				token:       tt.fields.token,
				Voucher:     tt.fields.Voucher,
			}
			got, got1, err := c.makeRequest(tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("makeRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("makeRequest() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("makeRequest() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestClient_urlBuilder(t *testing.T) {
	type fields struct {
		Credentials credentials.Credentials
		browser     *http.Client
		Url         *url.URL
		token       string
		Voucher     *voucher.Voucher
	}
	type args struct {
		endpoint string
		referer  string
	}
	var tests []struct {
		name   string
		fields fields
		args   args
		want   string
		want1  string
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				Credentials: tt.fields.Credentials,
				browser:     tt.fields.browser,
				Url:         tt.fields.Url,
				token:       tt.fields.token,
				Voucher:     tt.fields.Voucher,
			}
			got, got1 := c.urlBuilder(tt.args.endpoint, tt.args.referer)
			if got != tt.want {
				t.Errorf("urlBuilder() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("urlBuilder() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestRequestNewVoucherResponse_successful(t *testing.T) {
	type fields struct {
		Meta struct {
			Rc string `json:"rc"`
		}
		Data []struct {
			CreateTime int `json:"create_time"`
		}
	}
	var tests []struct {
		name   string
		fields fields
		want   bool
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := RequestNewVoucherResponse{
				Meta: tt.fields.Meta,
				Data: tt.fields.Data,
			}
			if got := v.successful(); got != tt.want {
				t.Errorf("successful() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_addBasicHeaders(t *testing.T) {
	type args struct {
		req *http.Request
	}
	var tests []struct {
		name string
		args args
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			addBasicHeaders(tt.args.req)
		})
	}
}

func Test_loggedIn(t *testing.T) {
	type args struct {
		responseBody string
	}
	var tests []struct {
		name string
		args args
		want bool
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := loggedIn(tt.args.responseBody); got != tt.want {
				t.Errorf("loggedIn() = %v, want %v", got, tt.want)
			}
		})
	}
}
