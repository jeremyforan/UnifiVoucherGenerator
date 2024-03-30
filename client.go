package UnifiVoucherGenerator

import (
	"github.com/jeremyforan/UnifiVoucherGenerator/credentials"
	"github.com/jeremyforan/UnifiVoucherGenerator/voucher"
	"log/slog"
	"net/http"
	"net/http/cookiejar"
	"net/url"
)

// todo: add additional logging
// todo: add ability to disable logging

// Client is the primary struct that interacts with the Unifi controller using http requests. It holds the credentials, http client, url, and token for the Unifi controller
// Before using the client, the Login method must be called to authenticate with the Unifi controller. If no error is returned, the client is ready to add vouchers.
type Client struct {
	Credentials credentials.Credentials
	browser     *http.Client
	Url         *url.URL
	token       string
	*voucher.Voucher
}

// NewClient creates a new Client struct to interact with the Unifi controller
func NewClient(username string, password string, url *url.URL) *Client {
	crd := credentials.NewCredentials(username, password)

	jar, _ := cookiejar.New(nil)

	return &Client{
		Credentials: crd,
		browser: &http.Client{
			Jar: jar,
		},
		Url: url,
	}
}

// Login sends a login request to the Unifi controller. If the login is successful, the csrf_token is stored in the client struct for future requests.
func (c *Client) Login() error {
	err := c.requestLogin()
	if err != nil {
		slog.Error("error logging in", "error", err)
		return err
	}
	return nil
}

// AddVoucher sends a request to the Unifi controller to add a voucher.
// If the request is successful, the voucher is stored in the client struct.
func (c *Client) AddVoucher(v *voucher.Voucher) error {
	c.Voucher = v

	err := c.requestAddVoucher()
	if err != nil {
		slog.Error("error adding voucher", "error", err)
		return err
	}
	c.Voucher.PublishedSuccesfully()

	vouchers, err := c.FetchVouchers()
	if err != nil {
		slog.Error("error fetching vouchers from Unifi", "error", err)
		return err
	}

	vUnifi, err := vouchers.getVoucherByID(c.Voucher.Id)
	if err != nil {
		slog.Error("error getting voucher from Unifi", "error", err)
		return err
	}

	ac, err := voucher.NewAccessCodeFromString(vUnifi.Code)
	if err != nil {
		slog.Error("error creating voucher code from string", "error", err)
		return err
	}

	c.Voucher.AC = ac
	return nil
}
