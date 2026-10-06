package UnifiVoucherGenerator

import (
	"errors"
	"testing"
)

func TestUnifiVouchers_getVoucherByID(t *testing.T) {
	list := UnifiVouchers{
		{Note: "one", Code: "1111122222"},
		{Note: "two", Code: "3333344444"},
	}

	got, err := list.getVoucherByID("two")
	if err != nil {
		t.Fatal(err)
	}
	if got.Code != "3333344444" {
		t.Errorf("code = %q", got.Code)
	}

	_, err = list.getVoucherByID("three")
	if !errors.Is(err, ErrVoucherNotFound) {
		t.Errorf("error = %v, want ErrVoucherNotFound", err)
	}

	_, err = UnifiVouchers(nil).getVoucherByID("one")
	if !errors.Is(err, ErrVoucherNotFound) {
		t.Errorf("empty list error = %v, want ErrVoucherNotFound", err)
	}
}
