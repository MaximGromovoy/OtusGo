package requests

type CalculateExchangeRequest struct {
	currencies     currencies
	fromAmount     float64
	commissionRate float64
}

func NewCalculateExchangeRequest(fromCurrency, toCurrency string, fromAmount, commissionRate float64) *CalculateExchangeRequest {
	return &CalculateExchangeRequest{
		currencies: currencies{
			fromCurrency: fromCurrency,
			toCurrency:   toCurrency,
		},
		fromAmount:     fromAmount,
		commissionRate: commissionRate,
	}
}

func (r *CalculateExchangeRequest) GetFromCurrency() string {
	return r.currencies.fromCurrency
}

func (r *CalculateExchangeRequest) GetToCurrency() string {
	return r.currencies.toCurrency
}

func (r *CalculateExchangeRequest) GetFromAmount() float64 {
	return r.fromAmount
}

func (r *CalculateExchangeRequest) GetCommissionRate() float64 {
	return r.commissionRate
}
