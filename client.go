package unifi

import (
	"io"
	"net/http"
	"net/http/cookiejar"
)

type Client struct {
	Credentials UnifiCredentials
	Client      *http.Client
	Url         string
	token       string
}

func NewClient(creds UnifiCredentials) *Client {
	jar, _ := cookiejar.New(nil)
	return &Client{
		Credentials: creds,
		Client: &http.Client{
			Jar: jar,
		},
		Url: unifiApiBaseUrl,
	}
}

func (c *Client) MakeRequest(req *http.Request) (string, []*http.Cookie, error) {
	res, err := c.Client.Do(req)
	if err != nil {
		return "", nil, err
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", nil, err
	}

	return string(body), res.Cookies(), nil
}
