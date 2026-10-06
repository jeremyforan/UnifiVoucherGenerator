package voucher

import (
	"encoding/json"
	"log/slog"
	"strings"
)

// Field names mirror the controller's create-voucher form: quota is the number of uses (0 = unlimited),
// expire_number/expire_unit give the lifetime (unit is minutes per step), up/down are Kbps and bytes is MB.

// Data is a struct that holds the parameters that can be selected online.
type Data struct {
	Quota            int    `json:"quota"`
	Note             string `json:"note"`
	NumberOfVouchers int    `json:"n"`
	ExpireNumber     string `json:"expire_number"`
	ExpireUnit       int    `json:"expire_unit"`
	Cmd              string `json:"cmd"`
	Up               int    `json:"up,omitempty"`
	Down             int    `json:"down,omitempty"`
	Bytes            string `json:"bytes,omitempty"`
}

func (v *Data) String() string {
	bytes, err := json.Marshal(v)
	if err != nil {
		slog.Error("Error marshaling to JSON", "error", err)
		return ""
	}
	return string(bytes)
}

// HttpPayload returns the NewVoucherRequestPayload struct as a strings.Reader to be used in
// as a http request body.
func (v *Data) HttpPayload() *strings.Reader {
	return strings.NewReader(v.String())
}

// LogGroup returns the NewVoucherRequestPayload struct as a slog.Value for logging.
func (v *Data) LogGroup() slog.Value {
	return slog.GroupValue(
		slog.Int("quota", v.Quota),
		slog.String("note", v.Note),
		slog.Int("amount", v.NumberOfVouchers),
		slog.String("expire_number", v.ExpireNumber),
		slog.Int("expire_unit", v.ExpireUnit),
		slog.String("cmd", v.Cmd),
	)
}
