package UnifiVoucherGenerator

import (
	"github.com/jeremyforan/UnifiVoucherGenerator/credentials"
	"github.com/jeremyforan/UnifiVoucherGenerator/voucher"
	"net/http"
	"net/url"
	"reflect"
	"testing"
)

func TestClient_requestAddVoucher(t *testing.T) {
	type fields struct {
		Credentials credentials.Credentials
		browser     *http.Client
		Url         *url.URL
		token       string
		Voucher     *voucher.Voucher
	}
	var tests []struct {
		name    string
		fields  fields
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
			if err := c.requestAddVoucher(); (err != nil) != tt.wantErr {
				t.Errorf("requestAddVoucher() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestClient_requestFetchPublishedVouchers(t *testing.T) {
	type fields struct {
		Credentials credentials.Credentials
		browser     *http.Client
		Url         *url.URL
		token       string
		Voucher     *voucher.Voucher
	}
	var tests []struct {
		name    string
		fields  fields
		want    UnifiVouchers
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
			got, err := c.requestFetchPublishedVouchers()
			if (err != nil) != tt.wantErr {
				t.Errorf("requestFetchPublishedVouchers() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("requestFetchPublishedVouchers() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClient_requestLogin(t *testing.T) {
	type fields struct {
		Credentials credentials.Credentials
		browser     *http.Client
		Url         *url.URL
		token       string
		Voucher     *voucher.Voucher
	}
	var tests []struct {
		name    string
		fields  fields
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
			if err := c.requestLogin(); (err != nil) != tt.wantErr {
				t.Errorf("requestLogin() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestClient_requestSelf(t *testing.T) {
	type fields struct {
		Credentials credentials.Credentials
		browser     *http.Client
		Url         *url.URL
		token       string
		Voucher     *voucher.Voucher
	}
	var tests []struct {
		name    string
		fields  fields
		want    string
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
			got, err := c.requestSelf()
			if (err != nil) != tt.wantErr {
				t.Errorf("requestSelf() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("requestSelf() got = %v, want %v", got, tt.want)
			}
		})
	}
}
