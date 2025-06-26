package handlers

import (
	"OtusGo/internal/interfaces"
	cs "OtusGo/internal/service/currencyService"
	csRequests "OtusGo/internal/service/currencyService/requests"
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"
)

type ExchangeHandler struct {
	currencyService       *cs.CurrencyService
	exchangeRatesProvider interfaces.ExchangeRatesProviderInterface
}

var commissionRate = 0.1 // Комиссия 10%

func NewExchangeHandler(currencyService *cs.CurrencyService, exchangeRatesProvider interfaces.ExchangeRatesProviderInterface) *ExchangeHandler {
	return &ExchangeHandler{
		currencyService:       currencyService,
		exchangeRatesProvider: exchangeRatesProvider,
	}
}

// ExchangeRequest HTTP запрос для обмена валют
type ExchangeRequestHTTP struct {
	UserID       int     `json:"user_id"`
	FromCurrency string  `json:"from_currency"`
	ToCurrency   string  `json:"to_currency"`
	FromAmount   float64 `json:"from_amount"`
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

type CalculateExchangeResponseHTTP struct {
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

// Exchange обрабатывает POST /api/exchange
func (h *ExchangeHandler) Exchange(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.writeErrorResponse(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var req ExchangeRequestHTTP
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeErrorResponse(w, http.StatusBadRequest, "invalid JSON: "+err.Error())
		return
	}

	// Преобразуем в внутренний формат
	exchangeReq := csRequests.NewExchangeRequest(req.UserID, req.FromCurrency, req.ToCurrency, req.FromAmount, commissionRate)

	// Выполняем обмен
	result, err := h.currencyService.Exchange(r.Context(), exchangeReq)
	if err != nil {
		slog.Error("Exchange failed", "error", err)
		h.writeErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	// Формируем успешный ответ
	response := ExchangeResponseHTTP{
		Success:        result.IsSuccessful,
		ToAmount:       result.Amount,
		ExchangeRate:   result.ExchangeRate,
		Commission:     result.Commission,
		CommissionRate: commissionRate,
	}

	h.writeJSONResponse(w, http.StatusOK, response)
}

// CalculateExchange обрабатывает GET /api/calculate?from=Dollar&to=Ruble&amount=100
func (h *ExchangeHandler) CalculateExchange(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.writeErrorResponse(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	fromCurrency := r.URL.Query().Get("from")
	toCurrency := r.URL.Query().Get("to")
	amountStr := r.URL.Query().Get("amount")

	if fromCurrency == "" || toCurrency == "" || amountStr == "" {
		h.writeErrorResponse(w, http.StatusBadRequest, "from, to currencies and amount are required")
		return
	}

	amount, err := strconv.ParseFloat(amountStr, 64)
	if err != nil {
		h.writeErrorResponse(w, http.StatusBadRequest, "invalid amount format")
		return
	}

	calculateExchangeRequest := csRequests.NewCalculateExchangeRequest(fromCurrency, toCurrency, amount, commissionRate)

	calculation, err := h.currencyService.CalculateExchangeAmount(r.Context(), calculateExchangeRequest)
	if err != nil {
		slog.Error("Failed to calculate exchange", "error", err)
		h.writeErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	response := CalculateExchangeResponseHTTP{
		Success:        calculation.IsSuccessful,
		ToAmount:       calculation.Amount,
		ExchangeRate:   calculation.ExchangeRate,
		Commission:     calculation.Commission,
		CommissionRate: commissionRate,
	}

	h.writeJSONResponse(w, http.StatusOK, response)
}

// GetRate обрабатывает GET /api/rate?from=Dollar&to=Ruble
func (h *ExchangeHandler) GetRate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.writeErrorResponse(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	fromCurrency := r.URL.Query().Get("from")
	toCurrency := r.URL.Query().Get("to")

	if fromCurrency == "" || toCurrency == "" {
		h.writeErrorResponse(w, http.StatusBadRequest, "from and to currencies are required")
		return
	}

	rate, err := h.exchangeRatesProvider.GetRate(r.Context(), fromCurrency, toCurrency)
	if err != nil {
		slog.Error("Failed to get exchange rate", "error", err)
		h.writeErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	response := GetRateResponseHTTP{
		Success:      true,
		FromCurrency: fromCurrency,
		ToCurrency:   toCurrency,
		ExchangeRate: rate,
	}

	h.writeJSONResponse(w, http.StatusOK, response)
}

func (h *ExchangeHandler) writeJSONResponse(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

func (h *ExchangeHandler) writeErrorResponse(w http.ResponseWriter, statusCode int, message string) {
	response := ExchangeResponseHTTP{
		Success: false,
		Error:   message,
	}
	h.writeJSONResponse(w, statusCode, response)
}
