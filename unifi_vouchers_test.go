package UnifiVoucherGenerator

import (
	"github.com/jeremyforan/UnifiVoucherGenerator/credentials"
	"github.com/jeremyforan/UnifiVoucherGenerator/voucher"
	"net/http"
	"net/url"
	"reflect"
	"testing"
)

func TestClient_FetchVouchers(t *testing.T) {
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
			got, err := c.FetchVouchers()
			if (err != nil) != tt.wantErr {
				t.Errorf("FetchVouchers() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FetchVouchers() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUnifiVouchers_getVoucherByID(t *testing.T) {
	type args struct {
		id string
	}
	var tests []struct {
		name    string
		v       UnifiVouchers
		args    args
		want    UnifiVoucher
		wantErr bool
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.v.getVoucherByID(tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("getVoucherByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getVoucherByID() got = %v, want %v", got, tt.want)
			}
		})
	}
}
