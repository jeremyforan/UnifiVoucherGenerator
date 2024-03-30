package voucher

type Expiry struct {
	Amount int
	Unit   ExpireUnit
}

func NewExpiry(amount int, unit ExpireUnit) Expiry {
	return Expiry{
		Amount: amount,
		Unit:   unit,
	}
}
