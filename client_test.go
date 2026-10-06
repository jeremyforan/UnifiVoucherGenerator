package UnifiVoucherGenerator

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/jeremyforan/UnifiVoucherGenerator/voucher"
)

// fakeController is a minimal stand-in for the UniFi Network Application's legacy API.
type fakeController struct {
	t            *testing.T
	loginOK      bool
	sendCSRF     bool
	createOK     bool
	listOK       bool
	code         string
	created      []map[string]any
	lastCSRF     string
	loginPayload string
}

func (f *fakeController) handler() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc(unifiApiLogin, func(w http.ResponseWriter, r *http.Request) {
		b, _ := io.ReadAll(r.Body)
		f.loginPayload = string(b)
		if !f.loginOK {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, `{"meta":{"rc":"error","msg":"api.err.Invalid"},"data":[]}`)
			return
		}
		if f.sendCSRF {
			http.SetCookie(w, &http.Cookie{Name: "csrf_token", Value: "tok123"})
		}
		http.SetCookie(w, &http.Cookie{Name: "unifises", Value: "session"})
		fmt.Fprint(w, `{"meta":{"rc":"ok"},"data":[]}`)
	})

	mux.HandleFunc(unifiApiCreateVoucher, func(w http.ResponseWriter, r *http.Request) {
		f.lastCSRF = r.Header.Get("X-Csrf-Token")
		var payload map[string]any
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			f.t.Errorf("create payload is not JSON: %v", err)
		}
		if !f.createOK {
			fmt.Fprint(w, `{"meta":{"rc":"error","msg":"api.err.InvalidPayload"},"data":[]}`)
			return
		}
		f.created = append(f.created, payload)
		fmt.Fprint(w, `{"meta":{"rc":"ok"},"data":[{"create_time":1722972000}]}`)
	})

	mux.HandleFunc(unifiApiVouchers, func(w http.ResponseWriter, r *http.Request) {
		if !f.listOK {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprint(w, `{"meta":{"rc":"error","msg":"api.err.LoginRequired"},"data":[]}`)
			return
		}
		list := VoucherListResponse{Meta: Meta{Rc: "ok"}}
		for _, c := range f.created {
			note, _ := c["note"].(string)
			list.Data = append(list.Data, UnifiVoucher{Note: note, Code: f.code, Quota: 1})
		}
		_ = json.NewEncoder(w).Encode(list)
	})

	return mux
}

func newTestClient(t *testing.T, f *fakeController) (*Client, *httptest.Server) {
	t.Helper()
	f.t = t
	ts := httptest.NewServer(f.handler())
	t.Cleanup(ts.Close)

	u, err := url.Parse(ts.URL)
	if err != nil {
		t.Fatal(err)
	}
	return NewClient("user@example.com", "p455w0rd", u), ts
}

func happyController() *fakeController {
	return &fakeController{loginOK: true, sendCSRF: true, createOK: true, listOK: true, code: "1234567890"}
}

func TestNewClient(t *testing.T) {
	u, _ := url.Parse("https://unifi.example.com:8443")
	c := NewClient("admin", "secret", u)

	if c.Credentials.Username != "admin" || c.Credentials.Password != "secret" {
		t.Errorf("credentials not stored: %+v", c.Credentials)
	}
	if c.Url != u {
		t.Errorf("url not stored")
	}
	if c.browser == nil || c.browser.Jar == nil {
		t.Fatalf("http client or cookie jar missing")
	}
	if c.browser.Timeout != DefaultTimeout {
		t.Errorf("timeout = %v, want %v", c.browser.Timeout, DefaultTimeout)
	}
	if c.Voucher != nil {
		t.Errorf("voucher should be nil before AddVoucher")
	}
}

func TestClient_StringBeforeAddVoucherDoesNotPanic(t *testing.T) {
	u, _ := url.Parse("https://unifi.example.com")
	c := NewClient("admin", "secret", u)
	if got := fmt.Sprintf("%v", c); got != "" {
		t.Errorf("String() on client without voucher = %q, want empty", got)
	}
}

func TestClient_SetHTTPClient(t *testing.T) {
	u, _ := url.Parse("https://unifi.example.com")
	c := NewClient("admin", "secret", u)
	original := c.browser

	c.SetHTTPClient(nil)
	if c.browser != original {
		t.Errorf("nil client should be ignored")
	}

	hc := &http.Client{}
	c.SetHTTPClient(hc)
	if c.browser != hc {
		t.Errorf("custom client not installed")
	}
	if hc.Jar == nil {
		t.Errorf("cookie jar should be added to a client without one")
	}
}

func TestClient_Login(t *testing.T) {
	t.Run("success stores csrf token", func(t *testing.T) {
		f := happyController()
		c, _ := newTestClient(t, f)
		if err := c.Login(); err != nil {
			t.Fatalf("Login() error = %v", err)
		}
		if c.token != "tok123" {
			t.Errorf("token = %q, want tok123", c.token)
		}
		if f.loginPayload != `{"username":"user@example.com","password":"p455w0rd","remember":true,"strict":true}` {
			t.Errorf("unexpected login payload %s", f.loginPayload)
		}
	})

	t.Run("rejected credentials", func(t *testing.T) {
		f := happyController()
		f.loginOK = false
		c, _ := newTestClient(t, f)
		err := c.Login()
		if !errors.Is(err, ErrLoginFailed) {
			t.Fatalf("Login() error = %v, want ErrLoginFailed", err)
		}
		if c.token != "" {
			t.Errorf("token should stay empty after failed login")
		}
	})

	t.Run("missing csrf cookie", func(t *testing.T) {
		f := happyController()
		f.sendCSRF = false
		c, _ := newTestClient(t, f)
		if err := c.Login(); !errors.Is(err, ErrCSRFTokenNotFound) {
			t.Fatalf("Login() error = %v, want ErrCSRFTokenNotFound", err)
		}
	})

	t.Run("unreachable controller returns transport error", func(t *testing.T) {
		f := happyController()
		c, ts := newTestClient(t, f)
		ts.Close()
		err := c.Login()
		if err == nil {
			t.Fatal("expected error")
		}
		if errors.Is(err, ErrLoginFailed) {
			t.Errorf("transport failure should not be reported as ErrLoginFailed: %v", err)
		}
	})
}

func TestClient_AddVoucher(t *testing.T) {
	t.Run("nil voucher", func(t *testing.T) {
		c, _ := newTestClient(t, happyController())
		if err := c.AddVoucher(nil); !errors.Is(err, ErrNilVoucher) {
			t.Errorf("error = %v, want ErrNilVoucher", err)
		}
	})

	t.Run("before login", func(t *testing.T) {
		c, _ := newTestClient(t, happyController())
		if err := c.AddVoucher(voucher.NewDefaultVoucher()); !errors.Is(err, ErrNotLoggedIn) {
			t.Errorf("error = %v, want ErrNotLoggedIn", err)
		}
	})

	t.Run("success populates access code", func(t *testing.T) {
		f := happyController()
		c, _ := newTestClient(t, f)
		if err := c.Login(); err != nil {
			t.Fatal(err)
		}

		v := voucher.NewMultiUseVoucher(5)
		v.SetExpire(2, voucher.Days)
		if err := c.AddVoucher(v); err != nil {
			t.Fatalf("AddVoucher() error = %v", err)
		}
		if !v.Published() {
			t.Errorf("voucher should be marked published")
		}
		if got := v.AccessCode().String(); got != "12345-67890" {
			t.Errorf("access code = %q, want 12345-67890", got)
		}
		if c.Voucher != v {
			t.Errorf("client should hold the voucher it just added")
		}
		if f.lastCSRF != "tok123" {
			t.Errorf("csrf header = %q, want tok123", f.lastCSRF)
		}
		if len(f.created) != 1 || f.created[0]["note"] != v.Id || f.created[0]["quota"] != float64(5) {
			t.Errorf("unexpected create payload %v", f.created)
		}
	})

	t.Run("custom id is used for lookup", func(t *testing.T) {
		f := happyController()
		c, _ := newTestClient(t, f)
		if err := c.Login(); err != nil {
			t.Fatal(err)
		}
		v := voucher.NewDefaultVoucher()
		v.SetId("front-door-button")
		if err := c.AddVoucher(v); err != nil {
			t.Fatalf("AddVoucher() with custom id error = %v", err)
		}
		if v.AccessCode().IsZero() {
			t.Errorf("access code not populated")
		}
	})

	t.Run("controller rejects voucher", func(t *testing.T) {
		f := happyController()
		f.createOK = false
		c, _ := newTestClient(t, f)
		if err := c.Login(); err != nil {
			t.Fatal(err)
		}
		v := voucher.NewDefaultVoucher()
		err := c.AddVoucher(v)
		if !errors.Is(err, ErrVoucherRequestFailed) {
			t.Fatalf("error = %v, want ErrVoucherRequestFailed", err)
		}
		if v.Published() {
			t.Errorf("rejected voucher must not be marked published")
		}
		if got := v.AccessCode().String(); got != "" {
			t.Errorf("access code should be empty, got %q", got)
		}
	})

	t.Run("expired session during fetch", func(t *testing.T) {
		f := happyController()
		f.listOK = false
		c, _ := newTestClient(t, f)
		if err := c.Login(); err != nil {
			t.Fatal(err)
		}
		err := c.AddVoucher(voucher.NewDefaultVoucher())
		if err == nil || errors.Is(err, ErrVoucherNotFound) {
			t.Fatalf("expected a session error rather than not-found, got %v", err)
		}
	})

	t.Run("malformed access code from controller", func(t *testing.T) {
		f := happyController()
		f.code = "ABCDE-12345"
		c, _ := newTestClient(t, f)
		if err := c.Login(); err != nil {
			t.Fatal(err)
		}
		if err := c.AddVoucher(voucher.NewDefaultVoucher()); err == nil {
			t.Fatal("expected error for malformed code")
		}
	})
}

func TestClient_FetchVouchers(t *testing.T) {
	t.Run("before login", func(t *testing.T) {
		c, _ := newTestClient(t, happyController())
		if _, err := c.FetchVouchers(); !errors.Is(err, ErrNotLoggedIn) {
			t.Errorf("error = %v, want ErrNotLoggedIn", err)
		}
	})

	t.Run("returns controller list", func(t *testing.T) {
		f := happyController()
		c, _ := newTestClient(t, f)
		if err := c.Login(); err != nil {
			t.Fatal(err)
		}
		for i := 0; i < 3; i++ {
			if err := c.AddVoucher(voucher.NewDefaultVoucher()); err != nil {
				t.Fatal(err)
			}
		}
		got, err := c.FetchVouchers()
		if err != nil {
			t.Fatal(err)
		}
		if len(got) != 3 {
			t.Errorf("len = %d, want 3", len(got))
		}
	})
}
