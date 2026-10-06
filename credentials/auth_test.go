package credentials

import (
	"encoding/json"
	"io"
	"log/slog"
	"strings"
	"testing"
)

func TestCredentials_String(t *testing.T) {
	c := NewCredentials("user@example.com", `pa"ss\word`)

	got := c.String()
	want := `{"username":"user@example.com","password":"pa\"ss\\word","remember":true,"strict":true}`
	if got != want {
		t.Errorf("String() = %s, want %s", got, want)
	}

	var back map[string]any
	if err := json.Unmarshal([]byte(got), &back); err != nil {
		t.Fatalf("payload is not valid JSON: %v", err)
	}
	if back["password"] != `pa"ss\word` {
		t.Errorf("password round trip = %v", back["password"])
	}
}

func TestCredentials_HttpPayload(t *testing.T) {
	c := NewCredentials("u", "p")
	b, err := io.ReadAll(c.HttpPayload())
	if err != nil {
		t.Fatal(err)
	}
	if string(b) != c.String() {
		t.Errorf("payload %s != String %s", b, c.String())
	}
}

func TestCredentials_LogValueHidesPassword(t *testing.T) {
	c := NewCredentials("u", "secret")
	var buf strings.Builder
	logger := slog.New(slog.NewTextHandler(&buf, nil))
	logger.Info("login", "creds", c)

	if strings.Contains(buf.String(), "secret") {
		t.Errorf("password leaked into log: %s", buf.String())
	}
	if !strings.Contains(buf.String(), "******") {
		t.Errorf("masked password missing: %s", buf.String())
	}
}
