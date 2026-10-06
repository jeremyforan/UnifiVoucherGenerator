package voucher

import "testing"

func TestNewAccessCodeFromString(t *testing.T) {
	tests := []struct {
		name    string
		in      string
		want    string
		wantErr bool
	}{
		{"digits only", "1234567890", "12345-67890", false},
		{"with dash", "12345-67890", "12345-67890", false},
		{"surrounding whitespace", " 12345-67890\n", "12345-67890", false},
		{"leading zeros", "00001-00002", "00001-00002", false},
		{"too short", "123456789", "", true},
		{"too long", "12345678901", "", true},
		{"letters", "ABCDE-12345", "", true},
		{"dash in wrong place", "1234-567890", "", true},
		{"empty", "", "", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewAccessCodeFromString(tt.in)
			if (err != nil) != tt.wantErr {
				t.Fatalf("error = %v, wantErr %v", err, tt.wantErr)
			}
			if got.String() != tt.want {
				t.Errorf("String() = %q, want %q", got.String(), tt.want)
			}
			if got.IsZero() != tt.wantErr {
				t.Errorf("IsZero() = %v, want %v", got.IsZero(), tt.wantErr)
			}
		})
	}
}

func TestVoucherCode(t *testing.T) {
	code := "1234567890"
	vc, err := NewAccessCodeFromString(code)
	if err != nil {
		t.Fatalf("Expected no error, got %s", err)
	}

	if len(vc.firstSet) != 5 || len(vc.secondSet) != 5 {
		t.Fatalf("Expected 5 and 5, got %d and %d", len(vc.firstSet), len(vc.secondSet))
	}

	wantFirst := []int{1, 2, 3, 4, 5}
	wantSecond := []int{6, 7, 8, 9, 0}
	for i := 0; i < 5; i++ {
		if vc.firstSet[i] != wantFirst[i] {
			t.Errorf("firstSet[%d] = %d, want %d", i, vc.firstSet[i], wantFirst[i])
		}
		if vc.secondSet[i] != wantSecond[i] {
			t.Errorf("secondSet[%d] = %d, want %d", i, vc.secondSet[i], wantSecond[i])
		}
	}
}

func TestAccessCode_ZeroValueString(t *testing.T) {
	var ac AccessCode
	if got := ac.String(); got != "" {
		t.Errorf("zero AccessCode String() = %q, want empty", got)
	}
	if !ac.IsZero() {
		t.Errorf("zero AccessCode should report IsZero")
	}
}
