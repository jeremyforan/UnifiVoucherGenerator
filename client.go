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

// Client is the primary struct that interacts with the Unifi controller using http requests
type Client struct {
	Credentials credentials.Credentials
	browser     *http.Client
	Url         *url.URL
	token       string
	*voucher.Voucher
}

// NewClient creates a new Client struct to interact with the Unifi controller
func NewClient(username string, password string, url *url.URL) *Client {
	credentials := credentials.NewCredentials(username, password)

	jar, _ := cookiejar.New(nil)

	return &Client{
		Credentials: credentials,
		browser: &http.Client{
			Jar: jar,
		},
		Url: url,
	}
}

func (c *Client) Login() error {
	err := c.requestLogin()
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) AddVoucher(v *voucher.Voucher) error {
	c.Voucher = v

	err := c.requestAddVoucher()
	if err != nil {
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
