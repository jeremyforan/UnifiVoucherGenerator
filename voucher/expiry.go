package voucher

type Expiry struct {
	Amount int
	Unit   VoucherExpireUnit
}

func NewExpiry(amount int, unit VoucherExpireUnit) Expiry {
	return Expiry{
		Amount: amount,
		Unit:   unit,
	}
}
