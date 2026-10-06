package voucher

import (
	"bytes"
	"github.com/google/uuid"
	"io"
	"log"
	"testing"
)

// TODO: write better tests, more rigid tests

// TestNewDefaultVoucher tests NewDefaultVoucher function
func TestNewDefaultVoucher(t *testing.T) {
	v := NewDefaultVoucher()

	v.SetId("87c1b126-86c2-4724-8e9b-a0933b1053a4")

	if v == nil {
		t.Errorf("NewDefaultVoucher returned nil")
	}
	if v.Id == "" || !isUUIDv4(v.Id) {
		t.Errorf("NewDefaultVoucher did not generate a valid UUID v4")
	}
	if v.data.ExpireUnit != int(Hours) {
		t.Errorf("NewDefaultVoucher expire unit is not set to Hours")
	}

	v.SetDataLimitMB(6)
	v.SetDownloadLimitMbps(8)
	v.SetUploadLimitMbps(7)
	v.SetExpire(25, Hours)

	payload := v.HttpPayload()

	expected := `{"quota":1,"note":"87c1b126-86c2-4724-8e9b-a0933b1053a4","n":1,"expire_number":"25","expire_unit":60,"cmd":"create-voucher","up":7000,"down":8000,"bytes":"6"}`

	if compareReaderToString(payload, expected) == false {
		t.Errorf("HttpPayload did not return the expected string")
	}

}

// TestNewSingleUseVoucher tests NewSingleUseVoucher function
func TestNewSingleUseVoucher(t *testing.T) {
	v := NewSingleUseVoucher()

	v.SetId("87c1b126-86c2-4724-8e9b-a0933b1053a4")

	if v == nil {
		t.Errorf("NewSingleUseVoucher returned nil")
	}
	if v.data.Quota != int(vSingleUse) {
		t.Errorf("NewSingleUseVoucher did not set Quota to single use")
	}

	v.SetDataLimitMB(6)
	v.SetDownloadLimitMbps(8)
	v.SetUploadLimitMbps(7)
	v.SetExpire(25, Hours)

	payload := v.HttpPayload()

	expected := `{"quota":1,"note":"87c1b126-86c2-4724-8e9b-a0933b1053a4","n":1,"expire_number":"25","expire_unit":60,"cmd":"create-voucher","up":7000,"down":8000,"bytes":"6"}`

	if compareReaderToString(payload, expected) == false {
		t.Errorf("HttpPayload did not return the expected string")
	}

}

// TestNewMultiUseVoucher tests NewMultiUseVoucher function with a specific quota
func TestNewMultiUseVoucher(t *testing.T) {
	quota := 10
	v := NewMultiUseVoucher(quota)

	v.SetId("87c1b126-86c2-4724-8e9b-a0933b1053a4")

	if v == nil {
		t.Errorf("NewMultiUseVoucher returned nil")
	}
	if v.data.Quota != quota {
		t.Errorf("NewMultiUseVoucher did not set the correct quota. Got %d, want %d", v.data.Quota, quota)
	}

	v.SetDataLimitMB(6)
	v.SetDownloadLimitMbps(8)
	v.SetUploadLimitMbps(7)
	v.SetExpire(25, Hours)

	payload := v.HttpPayload()

	expected := `{"quota":10,"note":"87c1b126-86c2-4724-8e9b-a0933b1053a4","n":1,"expire_number":"25","expire_unit":60,"cmd":"create-voucher","up":7000,"down":8000,"bytes":"6"}`

	if compareReaderToString(payload, expected) == false {
		t.Errorf("HttpPayload did not return the expected string")
	}

}

// TestNewUnlimitedUseVoucher tests NewUnlimitedUseVoucher function
func TestNewUnlimitedUseVoucher(t *testing.T) {
	v := NewUnlimitedUseVoucher()

	v.SetId("87c1b126-86c2-4724-8e9b-a0933b1053a4")

	if v == nil {
		t.Errorf("NewUnlimitedUseVoucher returned nil")
	}
	if v.data.Quota != int(vUnlimited) {
		t.Errorf("NewUnlimitedUseVoucher did not set Quota to unlimited")
	}

	v.SetDataLimitMB(6)
	v.SetDownloadLimitMbps(8)
	v.SetUploadLimitMbps(7)
	v.SetExpire(25, Hours)

	payload := v.HttpPayload()

	expected := `{"quota":0,"note":"87c1b126-86c2-4724-8e9b-a0933b1053a4","n":1,"expire_number":"25","expire_unit":60,"cmd":"create-voucher","up":7000,"down":8000,"bytes":"6"}`

	if compareReaderToString(payload, expected) == false {
		t.Errorf("HttpPayload did not return the expected string")
	}

}

func compareReaderToString(r *bytes.Reader, expected string) bool {
	// Read the contents of the *bytes.Reader
	content, err := io.ReadAll(r)
	if err != nil {
		log.Fatalf("Failed to read from bytes.Reader: %v", err)
	}

	// Convert the read content to a string
	contentStr := string(content)

	// Directly compare the strings
	return contentStr == expected
}

func TestVoucher_NilString(t *testing.T) {
	var v *Voucher
	if got := v.String(); got != "" {
		t.Errorf("nil Voucher String() = %q, want empty", got)
	}
}

func TestVoucher_SetIdKeepsIdAndNoteInSync(t *testing.T) {
	v := NewDefaultVoucher()
	v.SetId("front-door")
	if v.Id != "front-door" {
		t.Errorf("Id = %q, want front-door", v.Id)
	}
	if v.data.Note != "front-door" {
		t.Errorf("Note = %q, want front-door", v.data.Note)
	}
}

func TestVoucher_DefaultExpiryIs24Hours(t *testing.T) {
	for name, v := range map[string]*Voucher{
		"default":   NewDefaultVoucher(),
		"single":    NewSingleUseVoucher(),
		"multi":     NewMultiUseVoucher(3),
		"unlimited": NewUnlimitedUseVoucher(),
	} {
		if v.data.ExpireNumber != "24" || v.data.ExpireUnit != int(Hours) {
			t.Errorf("%s: expiry = %s unit %d, want 24 hours", name, v.data.ExpireNumber, v.data.ExpireUnit)
		}
		if v.data.Note != v.Id {
			t.Errorf("%s: note %q does not match id %q", name, v.data.Note, v.Id)
		}
		if v.Published() {
			t.Errorf("%s: new voucher should not be published", name)
		}
	}
}

func TestVoucher_DefaultHasNoLimits(t *testing.T) {
	v := NewDefaultVoucher()
	v.SetId("fixed")
	expected := `{"quota":1,"note":"fixed","n":1,"expire_number":"24","expire_unit":60,"cmd":"create-voucher"}`
	if !compareReaderToString(v.HttpPayload(), expected) {
		t.Errorf("payload = %s, want %s", v.String(), expected)
	}
}

func isUUIDv4(s string) bool {
	u, err := uuid.Parse(s)
	return err == nil && u.Version() == 4
}
