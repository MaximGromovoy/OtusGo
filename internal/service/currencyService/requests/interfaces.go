package requests

// Описывает наличие пользователя в запросе
type UserRequest interface {
	GetUserId() int
}

// Описывает наличие валют в запросе
type CurrenciesRequest interface {
	GetFromCurrency() string
	GetToCurrency() string
}

// Описывает наличие суммы в запросе
type AmountRequest interface {
	GetFromAmount() float64
}

// Описывает наличие комиссии в запросе
type CommissionRateRequest interface {
	GetCommissionRate() float64
}
