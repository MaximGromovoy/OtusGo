package operations

import "math"

type ExchangeOperationRequest struct {
	FromAmount     float64 // Сумма для обмена
	ExchangeRate   float64 // Курс обмена
	CommissionRate float64 // Процент комиссии
}

type ExchangeOperationResult struct {
	ToAmountBeforeCommission float64
	ToAmountAfterCommission  float64
	Commission               float64
}

// Exchange выполняет расчёт обмена
func Exchange(req *ExchangeOperationRequest) ExchangeOperationResult {
	toAmountBeforeCommission := req.FromAmount * req.ExchangeRate
	commission := toAmountBeforeCommission * (req.CommissionRate / 100.0)
	commission = math.Round(commission*100) / 100
	toAmountAfterCommission := toAmountBeforeCommission - commission

	return ExchangeOperationResult{
		ToAmountBeforeCommission: toAmountBeforeCommission,
		ToAmountAfterCommission:  toAmountAfterCommission,
		Commission:               commission,
	}
}
