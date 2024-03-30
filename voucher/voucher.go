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
	AC        AccessCode
}

// NewDefaultVoucher creates a new Voucher struct to be used to create a new voucher.
func NewDefaultVoucher() *Voucher {
	id := uuid.NewV4().String()
	return &Voucher{
		Id:        id,
		published: false,
		data: Data{
			Note:             id,
			Quota:            1,
			NumberOfVouchers: 1,
			ExpireNumber:     24,
			ExpireUnit:       int(vHours),
			Cmd:              createVoucher,
		},
	}
}

// NewSingleUseVoucher creates a new Single Use Voucher.
func NewSingleUseVoucher() *Voucher {
	v := blankVoucher()
	v.data = Data{
		Quota: int(vSingleUse),
		Cmd:   createVoucher,
	}
	return &v
}

// NewMultiUseVoucher creates a new Multi Use Voucher.
func NewMultiUseVoucher(quota int) *Voucher {
	v := blankVoucher()
	v.data = Data{
		Quota: quota,
		Cmd:   createVoucher,
	}
	return &v
}

// NewUnlimitedUseVoucher creates a new Unlimited use Voucher.
func NewUnlimitedUseVoucher() *Voucher {
	v := blankVoucher()
	v.data = Data{
		Quota: int(vUnlimited),
		Cmd:   createVoucher,
	}
	return &v
}

func (v *Voucher) String() string {
	return v.data.String()
}

func (v *Voucher) HttpPayload() *strings.Reader {
	return v.data.HttpPayload()
}

func (v *Voucher) AccessCode() AccessCode {
	return v.AC
}

func (v *Voucher) PublishedSuccesfully() {
	v.published = true
}

//todo: add download limit
//todo: add upload limit
//todo: add bandwidth limit
