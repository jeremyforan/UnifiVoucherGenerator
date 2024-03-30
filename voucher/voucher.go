package voucher

import (
	"github.com/satori/go.uuid"
	"strings"
)

//todo: I think this should tweak the names a little bit.

// Voucher is a struct that holds the information needed to create a new voucher.
type Voucher struct {
	Id        string
	published bool
	data      Data
	Code
}

// NewDefaultVoucher creates a new Voucher struct to be used to create a new voucher.
func NewDefaultVoucher() Voucher {
	id := uuid.NewV4().String()
	return Voucher{
		Id:        id,
		published: false,
		data: Data{
			Note:             id,
			Quota:            1,
			NumberOfVouchers: 1,
			ExpireNumber:     24,
			ExpireUnit:       int(vHours),
		},
	}
}

// NewSingleUseVoucher creates a new Single Use Voucher.
func NewSingleUseVoucher() Voucher {
	v := blankVoucher()
	v.data = Data{
		Quota: int(vSingleUse),
		Cmd:   string(createVoucher),
	}
	return v
}

// NewMultiUseVoucher creates a new Multi Use Voucher.
func NewMultiUseVoucher(quota int) Voucher {
	v := blankVoucher()
	v.data = Data{
		Quota: quota,
		Cmd:   string(createVoucher),
	}
	return v
}

// NewUnlimitedUseVoucher creates a new Unlimited use Voucher.
func NewUnlimitedUseVoucher() Voucher {
	v := blankVoucher()
	v.data = Data{
		Quota: int(vUnlimited),
		Cmd:   string(createVoucher),
	}
	return v
}

func (v *Voucher) String() string {
	return v.data.String()
}

func (v *Voucher) HttpPayload() *strings.Reader {
	return v.data.HttpPayload()
}

func (v *Voucher) AccessCode() Code {
	return v.Code
}

func (v *Voucher) PublishedSuccesfully() {
	v.published = true
}
