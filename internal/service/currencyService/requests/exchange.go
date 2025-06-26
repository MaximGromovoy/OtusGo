package requests

type ExchangeRequest struct {
	userId         int
	currencies     currencies
	fromAmount     float64
	commissionRate float64
}

func NewExchangeRequest(userId int, fromCurrency, toCurrency string, fromAmount, commissionRate float64) *ExchangeRequest {
	return &ExchangeRequest{
		userId: userId,
		currencies: currencies{
			fromCurrency: fromCurrency,
			toCurrency:   toCurrency,
		},
		fromAmount:     fromAmount,
		commissionRate: commissionRate,
	}
}

func (r *ExchangeRequest) GetUserId() int {
	return r.userId
}

func (r *ExchangeRequest) GetFromCurrency() string {
	return r.currencies.fromCurrency
}

func (r *ExchangeRequest) GetToCurrency() string {
	return r.currencies.toCurrency
}

func (r *ExchangeRequest) GetFromAmount() float64 {
	return r.fromAmount
}

func (r *ExchangeRequest) GetCommissionRate() float64 {
	return r.commissionRate
}
