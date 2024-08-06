package UnifiVoucherGenerator

import (
	"fmt"
	"github.com/jeremyforan/UnifiVoucherGenerator/credentials"
	"github.com/jeremyforan/UnifiVoucherGenerator/voucher"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"testing"
)

func TestClient_AddVoucher(t *testing.T) {
	ts := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprintln(w, "Hello, client")
			}))

	defer ts.Close()

	type fields struct {
		Credentials credentials.Credentials
		browser     *http.Client
		Url         *url.URL
		token       string
		Voucher     *voucher.Voucher
	}

	type args struct {
		v *voucher.Voucher
	}

	var tests []struct {
		name    string
		fields  fields
		args    args
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
			if err := c.AddVoucher(tt.args.v); (err != nil) != tt.wantErr {
				t.Errorf("AddVoucher() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestClient_Login(t *testing.T) {
	ts := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprintln(w, "Hello, client")
			}))

	defer ts.Close()

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
			if err := c.Login(); (err != nil) != tt.wantErr {
				t.Errorf("Login() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNewClient(t *testing.T) {
	validURL, _ := url.Parse("https://unifi.jeremyforan.com")
	//invalidURL, _ := url.Parse("https://unifi.jeremyforan.com")

	validUsername := "4dm1n157r470r"
	validPassword := "p455w0rd"

	//validUsername := "administrator"
	//validPassword := "password"

	type args struct {
		username string
		password string
		url      *url.URL
	}
	var tests []struct {
		name string
		args args
		want *Client
	}

	tests = append(tests, struct {
		name string
		args args
		want *Client
	}{
		name: "New Client - valid username, valid password, valid URL",
		args: args{
			username: validUsername,
			password: validPassword,
			url:      validURL,
		},
		want: &Client{
			Credentials: credentials.Credentials{
				Username: validUsername,
				Password: validPassword,
			},
		},
	})

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fmt.Println(tt.args)
			if got := NewClient(tt.args.username, tt.args.password, tt.args.url); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewClient() = %v, want %v", got, tt.want)
			}
		})
	}
}

//
//base, err := url.Parse("https://unifi.jeremyforan.com")
//if err != nil {
//	slog.Error("failed to parse Unifi controller URL")
//	panic(err)
//}
//
//slog.Info("creating Unifi client", "url", base)
//client := UnifiVoucherGenerator.NewClient("jeremy.foran@gmail.com", "5):F02~2p>4_>xBw+V26*+QeYUvug9t%", base)
//
//err = client.Login()
//if err != nil {
//slog.Error("failed to login to Unifi controller")
//panic(err)
//}
//
//slog.Info("adding voucher")
//v := voucher.NewMultiUseVoucher(5)
//v.SetExpire(96, voucher.Hours)
//
//slog.Info("adding voucher to unifi controller")
//err = client.AddVoucher(v)
//if err != nil {
//slog.Error("unable to add voucher")
//panic(err)
//}
//slog.Info("voucher added", "voucher", v)
//
//slog.Info("printing voucher", "voucher", v.AccessCode().String())
//err = s.PrintVoucher(v.AccessCode().String())
//if err != nil {
//slog.Error("failed to print voucher")
//panic(err)
//}
