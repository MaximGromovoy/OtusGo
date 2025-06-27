package handlers

// ExchangeRequest HTTP запрос для обмена валют
type ExchangeRequestHTTP struct {
	UserID         int     `json:"user_id"`
	FromCurrency   string  `json:"from_currency"`
	ToCurrency     string  `json:"to_currency"`
	FromAmount     float64 `json:"from_amount"`
	CommissionRate float64 `json:"commission_rate"`
}

// ExchangeResponse HTTP ответ с результатом обмена
type ExchangeResponseHTTP struct {
	Success        bool    `json:"success"`
	ToAmount       float64 `json:"to_amount,omitempty"`
	ExchangeRate   float64 `json:"exchange_rate,omitempty"`
	Commission     float64 `json:"commission,omitempty"`
	CommissionRate float64 `json:"commission_rate,omitempty"`
	Error          string  `json:"error,omitempty"`
}

type GetRateResponseHTTP struct {
	Success      bool    `json:"success"`
	FromCurrency string  `json:"from_currency,omitempty"`
	ToCurrency   string  `json:"to_currency,omitempty"`
	ExchangeRate float64 `json:"exchange_rate,omitempty"`
	Error        string  `json:"error,omitempty"`
}
