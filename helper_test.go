package UnifiVoucherGenerator

import (
	"net/http"
	"net/url"
	"testing"
)

func TestClient_urlBuilder(t *testing.T) {
	tests := []struct {
		name     string
		base     string
		endpoint string
		referer  string
		want     string
		want1    string
	}{
		{
			name:     "plain host",
			base:     "https://unifi.example.com:8443",
			endpoint: unifiApiLogin,
			referer:  unifiApiLoginReferer,
			want:     "https://unifi.example.com:8443/api/login",
			want1:    "https://unifi.example.com:8443/manage/account/login",
		},
		{
			name:     "trailing slash is not doubled",
			base:     "https://unifi.example.com/",
			endpoint: unifiApiVouchers,
			referer:  unifiApiVoucherReferer,
			want:     "https://unifi.example.com/api/s/default/stat/voucher",
			want1:    "https://unifi.example.com/manage/default/hotspot/vouchers",
		},
		{
			name:     "base path is preserved",
			base:     "https://example.com/network",
			endpoint: unifiApiCreateVoucher,
			referer:  unifiApiVoucherReferer,
			want:     "https://example.com/network/api/s/default/cmd/hotspot",
			want1:    "https://example.com/network/manage/default/hotspot/vouchers",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u, err := url.Parse(tt.base)
			if err != nil {
				t.Fatal(err)
			}
			c := &Client{Url: u}
			got, got1 := c.urlBuilder(tt.endpoint, tt.referer)
			if got != tt.want {
				t.Errorf("endpoint = %q, want %q", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("referer = %q, want %q", got1, tt.want1)
			}
		})
	}

	t.Run("nil url", func(t *testing.T) {
		c := &Client{}
		got, got1 := c.urlBuilder(unifiApiLogin, unifiApiLoginReferer)
		if got != "" || got1 != "" {
			t.Errorf("expected empty urls for nil base, got %q %q", got, got1)
		}
	})
}

func TestClient_urlHelpers(t *testing.T) {
	u, _ := url.Parse("https://unifi.example.com")
	c := &Client{Url: u}

	if a, _ := c.loginUrls(); a != "https://unifi.example.com/api/login" {
		t.Errorf("loginUrls = %q", a)
	}
	if a, _ := c.addVoucherUrls(); a != "https://unifi.example.com/api/s/default/cmd/hotspot" {
		t.Errorf("addVoucherUrls = %q", a)
	}
	if a, _ := c.fetchVouchersUrl(); a != "https://unifi.example.com/api/s/default/stat/voucher" {
		t.Errorf("fetchVouchersUrl = %q", a)
	}
}

func Test_addBasicHeaders(t *testing.T) {
	req, _ := http.NewRequest(http.MethodGet, "https://example.com", nil)
	addBasicHeaders(req)

	want := map[string]string{
		"Accept":        "*/*",
		"Content-Type":  "application/json; charset=utf-8",
		"Cache-Control": "no-cache",
		"DNT":           "1",
	}
	for k, v := range want {
		if got := req.Header.Get(k); got != v {
			t.Errorf("header %s = %q, want %q", k, got, v)
		}
	}
}

func Test_loggedIn(t *testing.T) {
	tests := []struct {
		name string
		body string
		want bool
	}{
		{"ok", `{"meta":{"rc":"ok"},"data":[]}`, true},
		{"error", `{"meta":{"rc":"error","msg":"api.err.Invalid"},"data":[]}`, false},
		{"empty", ``, false},
		{"html", `<html>login</html>`, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := loggedIn(tt.body); got != tt.want {
				t.Errorf("loggedIn() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_csrfTokenFromCookies(t *testing.T) {
	tests := []struct {
		name    string
		cookies []*http.Cookie
		want    string
		wantOK  bool
	}{
		{"present", []*http.Cookie{{Name: "unifises", Value: "s"}, {Name: "csrf_token", Value: "abc"}}, "abc", true},
		{"absent", []*http.Cookie{{Name: "unifises", Value: "s"}}, "", false},
		{"empty value", []*http.Cookie{{Name: "csrf_token", Value: ""}}, "", false},
		{"no cookies", nil, "", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := csrfTokenFromCookies(tt.cookies)
			if got != tt.want || ok != tt.wantOK {
				t.Errorf("got (%q, %v), want (%q, %v)", got, ok, tt.want, tt.wantOK)
			}
		})
	}
}

func TestRequestNewVoucherResponse_successful(t *testing.T) {
	var ok, bad RequestNewVoucherResponse
	ok.Meta.Rc = "ok"
	bad.Meta.Rc = "error"

	if !ok.successful() {
		t.Errorf("rc ok should be successful")
	}
	if bad.successful() {
		t.Errorf("rc error should not be successful")
	}
}
