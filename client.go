package UnifiVoucherGenerator

import (
	"github.com/jeremyforan/UnifiVoucherGenerator/credentials"
	"github.com/jeremyforan/UnifiVoucherGenerator/voucher"
	"net/http"
	"net/http/cookiejar"
	"net/url"
)

// todo: add additional logging

// Client is the primary struct that interacts with the Unifi controller using http requests
type Client struct {
	Credentials credentials.Credentials
	client      *http.Client
	Url         url.URL
	token       string
}

// NewClient creates a new Client struct to interact with the Unifi controller
func NewClient(credentials credentials.Credentials, url url.URL) *Client {
	jar, _ := cookiejar.New(nil)
	return &Client{
		Credentials: credentials,
		client: &http.Client{
			Jar: jar,
		},
		Url: url,
	}
}

func (c *Client) Login() error {

	return nil
}

func (c *Client) AddVoucher(v voucher.Voucher) error {
	return nil
}
