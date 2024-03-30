package voucher

import "testing"

func TestVoucherCode(t *testing.T) {
	code := "1234567890"
	vc, err := NewAccessCodeFromString(code)
	if err != nil {
		t.Errorf("Expected no error, got %s", err)
	}

	if len(vc.firstSet) != 5 {
		t.Errorf("Expected 5, got %d", len(vc.firstSet))
	}

	if len(vc.secondSet) != 5 {
		t.Errorf("Expected 5, got %d", len(vc.secondSet))
	}

	for i := 0; i < 5; i++ {
		if vc.firstSet[i] != i+1 {
			t.Errorf("Expected %d, got %d", int(code[i]), vc.firstSet[i])
		}
	}

	for i := 0; i < 5; i++ {
		if vc.secondSet[i] != i+1 {
			t.Errorf("Expected %d, got %d", int(code[i]), vc.firstSet[i])
		}
	}
}
