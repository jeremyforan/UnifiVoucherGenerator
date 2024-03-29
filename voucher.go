package unifi

import (
	"fmt"
	"strings"
)

type NewVoucherRequestPayload struct {
	Quota        int    `json:"quota"`
	Note         string `json:"note"`
	N            int    `json:"n"`
	ExpireNumber int    `json:"expire_number"`
	ExpireUnit   int    `json:"expire_unit"`
	Cmd          string `json:"cmd"`
}

func NewVoucherPayload() NewVoucherRequestPayload {
	return NewVoucherRequestPayload{
		Quota:        1,
		Note:         NoteTimeStamp(),
		N:            1,
		ExpireNumber: 24,
		ExpireUnit:   60,
		Cmd:          "create-voucher",
	}
}

func (v *NewVoucherRequestPayload) String() string {
	return fmt.Sprintf(`{"quota":%d,"note":"%s","n":%d,"expire_number":%d,"expire_unit":%d,"cmd":"%s"}`, v.Quota, v.Note, v.N, v.ExpireNumber, v.ExpireUnit, v.Cmd)
}

func (v *NewVoucherRequestPayload) HttpPayload() *strings.Reader {
	return strings.NewReader(v.String())
}
