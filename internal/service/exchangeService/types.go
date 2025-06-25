package exchangeService

import (
	"OtusGo/internal/transaction"
)

// ExchangeRequest запрос на обмен валют
type ExchangeRequest struct {
	UserID       int     // ID пользователя
	FromCurrency string  // Исходная валюта (Dollar, Euro, Ruble, Lira)
	ToCurrency   string  // Целевая валюта
	FromAmount   float64 // Сумма для обмена
}

// ExchangeResult результат обмена валют
type ExchangeResult struct {
	Transaction    *transaction.Transaction // Созданная транзакция
	ToAmount       float64                  // Итоговая сумма после обмена
	ExchangeRate   float64                  // Использованный курс
	Commission     float64                  // Размер комиссии
	CommissionRate float64                  // Процент комиссии
}

// ExchangeCalculation результат расчета обмена
type ExchangeCalculation struct {
	FromCurrency             string  // Исходная валюта
	ToCurrency               string  // Целевая валюта
	FromAmount               float64 // Исходная сумма
	ToAmountBeforeCommission float64 // Сумма до вычета комиссии
	ToAmountAfterCommission  float64 // Сумма после вычета комиссии
	ExchangeRate             float64 // Курс обмена
	Commission               float64 // Размер комиссии
	CommissionRate           float64 // Процент комиссии
}
