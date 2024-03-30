package voucher

// voucherExpireUnit is the unit of time that the voucher will expire in
type VoucherExpireUnit int

const (
	vMinutes VoucherExpireUnit = 1
	vHours   VoucherExpireUnit = 60
	vDays    VoucherExpireUnit = 1440
)

type VoucherUsageType int

const (
	vSingleUse VoucherUsageType = 1
	vUnlimited VoucherUsageType = 0

	//	vMultiUse is set during constructor
)
