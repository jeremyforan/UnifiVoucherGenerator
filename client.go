package UnifiVoucherGenerator

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"time"

	"github.com/jeremyforan/UnifiVoucherGenerator/credentials"
	"github.com/jeremyforan/UnifiVoucherGenerator/voucher"
)

// DefaultTimeout is the request timeout applied to the http.Client created by NewClient.
const DefaultTimeout = 30 * time.Second

// Sentinel errors returned by the client. They are wrapped with additional context, so
// compare them with errors.Is.
var (
	ErrLoginFailed          = errors.New("login failed")
	ErrCSRFTokenNotFound    = errors.New("csrf_token not found")
	ErrNotLoggedIn          = errors.New("not logged in")
	ErrNilVoucher           = errors.New("voucher is nil")
	ErrVoucherRequestFailed = errors.New("voucher request failed")
	ErrVoucherNotFound      = errors.New("voucher not found")
)

// Client is the primary struct that interacts with the Unifi controller using http requests. It holds the credentials, http client, url, and token for the Unifi controller
// Before using the client, the Login method must be called to authenticate with the Unifi controller. If no error is returned, the client is ready to add vouchers.
//
// A Client is not safe for concurrent use. The embedded Voucher is the voucher most
// recently passed to AddVoucher and is nil until then.
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
			Jar:     jar,
			Timeout: DefaultTimeout,
		},
		Url: url,
	}
}

// SetHTTPClient replaces the http.Client used to talk to the controller. Use this to set a
// custom timeout or TLS configuration, for example to trust a controller's self-signed
// certificate. If the supplied client has no cookie jar one is added, because the controller
// session is cookie based. A nil client is ignored.
func (c *Client) SetHTTPClient(hc *http.Client) {
	if hc == nil {
		return
	}
	if hc.Jar == nil {
		hc.Jar, _ = cookiejar.New(nil)
	}
	c.browser = hc
}

// Login sends a login request to the Unifi controller. If the login is successful, the csrf_token is stored in the client struct for future requests.
func (c *Client) Login() error {
	if err := c.requestLogin(); err != nil {
		slog.Error("error logging in", "error", err)
		return err
	}
	return nil
}

// AddVoucher sends a request to the Unifi controller to add a voucher.
// If the request is successful, the voucher is stored in the client struct and its
// access code is populated, so v.AccessCode() can be used afterwards.
func (c *Client) AddVoucher(v *voucher.Voucher) error {
	if v == nil {
		return ErrNilVoucher
	}
	if c.token == "" {
		return ErrNotLoggedIn
	}

	c.Voucher = v

	if err := c.requestAddVoucher(); err != nil {
		slog.Error("error adding voucher", "error", err)
		return err
	}
	v.PublishedSuccesfully()

	vouchers, err := c.FetchVouchers()
	if err != nil {
		slog.Error("error fetching vouchers from Unifi", "error", err)
		return err
	}

	vUnifi, err := vouchers.getVoucherByID(v.Id)
	if err != nil {
		slog.Error("error getting voucher from Unifi", "error", err)
		return err
	}

	ac, err := voucher.NewAccessCodeFromString(vUnifi.Code)
	if err != nil {
		err = fmt.Errorf("voucher %q returned by controller: %w", vUnifi.Code, err)
		slog.Error("error creating voucher code from string", "error", err)
		return err
	}

	v.AC = ac
	return nil
}
