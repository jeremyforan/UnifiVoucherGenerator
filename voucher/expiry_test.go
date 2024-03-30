package voucher

import "testing"

//todo: need to run better tests

func TestNewExpiry(t *testing.T) {
	amount := 24
	unit := vHours
	expiry := NewExpiry(amount, unit)

	if expiry.Amount != amount {
		t.Errorf("Expected %d, got %d", amount, expiry.Amount)
	}

	if expiry.Unit != unit {
		t.Errorf("Expected %d, got %d", unit, expiry.Unit)
	}
}
