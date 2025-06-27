package exchangeService

type ExchangeRequest struct {
	Amount         float64
	CommissionRate float64
	ExchangeRate   float64
}

func NewExchangeRequest(amount, commissionRate, exchangeRate float64) *ExchangeRequest {
	return &ExchangeRequest{
		Amount:         amount,
		CommissionRate: commissionRate,
		ExchangeRate:   exchangeRate,
	}
}

type ExchangeResponse struct {
	Amount     float64
	Commission float64
}

func NewExchangeResponse(amount, commission float64) *ExchangeResponse {
	return &ExchangeResponse{
		Amount:     amount,
		Commission: commission,
	}
}
