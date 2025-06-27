package exchangeService

import (
	"math"
)

type ExchangeService struct{}

func NewExchangeService() *ExchangeService {
	return &ExchangeService{}
}

// Exchange выполняет обмен валют
func (s *ExchangeService) Exchange(req *ExchangeRequest) *ExchangeResponse {
	toAmountBeforeCommission := req.Amount * req.ExchangeRate
	commission := toAmountBeforeCommission * (req.CommissionRate / 100.0)
	commission = math.Round(commission*100) / 100
	toAmountAfterCommission := toAmountBeforeCommission - commission

	return NewExchangeResponse(toAmountAfterCommission, commission)
}
