package handlers

import (
	"OtusGo/internal/service/exchangeService"
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"
)

type ExchangeHandler struct {
	exchangeService *exchangeService.ExchangeService
}

func NewExchangeHandler(exchangeService *exchangeService.ExchangeService) *ExchangeHandler {
	return &ExchangeHandler{
		exchangeService: exchangeService,
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
	TransactionID  int     `json:"transaction_id,omitempty"`
	ToAmount       float64 `json:"to_amount,omitempty"`
	ExchangeRate   float64 `json:"exchange_rate,omitempty"`
	Commission     float64 `json:"commission,omitempty"`
	CommissionRate float64 `json:"commission_rate,omitempty"`
	Error          string  `json:"error,omitempty"`
}

// ExchangeCurrency обрабатывает POST /api/exchange
func (h *ExchangeHandler) ExchangeCurrency(w http.ResponseWriter, r *http.Request) {
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
	exchangeReq := &exchangeService.ExchangeRequest{
		UserID:       req.UserID,
		FromCurrency: req.FromCurrency,
		ToCurrency:   req.ToCurrency,
		FromAmount:   req.FromAmount,
	}

	// Выполняем обмен
	result, err := h.exchangeService.ExchangeCurrency(r.Context(), exchangeReq)
	if err != nil {
		slog.Error("Exchange failed", "error", err)
		h.writeErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	// Формируем успешный ответ
	response := ExchangeResponseHTTP{
		Success:        true,
		TransactionID:  result.Transaction.ID,
		ToAmount:       result.ToAmount,
		ExchangeRate:   result.ExchangeRate,
		Commission:     result.Commission,
		CommissionRate: result.CommissionRate,
	}

	h.writeJSONResponse(w, http.StatusOK, response)
}

// GetExchangeRate обрабатывает GET /api/rates?from=USD&to=RUB
func (h *ExchangeHandler) GetExchangeRate(w http.ResponseWriter, r *http.Request) {
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

	rate, err := h.exchangeService.GetExchangeRate(r.Context(), fromCurrency, toCurrency)
	if err != nil {
		slog.Error("Failed to get exchange rate", "error", err)
		h.writeErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	response := map[string]interface{}{
		"success":       true,
		"from_currency": fromCurrency,
		"to_currency":   toCurrency,
		"exchange_rate": rate,
	}

	h.writeJSONResponse(w, http.StatusOK, response)
}

// CalculateExchange обрабатывает GET /api/calculate?from=USD&to=RUB&amount=100
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

	calculation, err := h.exchangeService.CalculateExchangeAmount(r.Context(), fromCurrency, toCurrency, amount)
	if err != nil {
		slog.Error("Failed to calculate exchange", "error", err)
		h.writeErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	response := map[string]interface{}{
		"success":                     true,
		"from_currency":               calculation.FromCurrency,
		"to_currency":                 calculation.ToCurrency,
		"from_amount":                 calculation.FromAmount,
		"to_amount_before_commission": calculation.ToAmountBeforeCommission,
		"to_amount_after_commission":  calculation.ToAmountAfterCommission,
		"exchange_rate":               calculation.ExchangeRate,
		"commission":                  calculation.Commission,
		"commission_rate":             calculation.CommissionRate,
	}

	h.writeJSONResponse(w, http.StatusOK, response)
}

// RefreshRates обрабатывает POST /api/refresh-rates
func (h *ExchangeHandler) RefreshRates(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.writeErrorResponse(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	err := h.exchangeService.RefreshRates(r.Context())
	if err != nil {
		slog.Error("Failed to refresh rates", "error", err)
		h.writeErrorResponse(w, http.StatusInternalServerError, "failed to refresh rates: "+err.Error())
		return
	}

	response := map[string]interface{}{
		"success": true,
		"message": "rates refreshed successfully",
	}

	h.writeJSONResponse(w, http.StatusOK, response)
}

// GetCachedRates обрабатывает GET /api/cached-rates
func (h *ExchangeHandler) GetCachedRates(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.writeErrorResponse(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	rates, err := h.exchangeService.GetCachedRates(r.Context())
	if err != nil {
		h.writeErrorResponse(w, http.StatusInternalServerError, "failed to get cached rates: "+err.Error())
		return
	}
	if rates == nil {
		h.writeErrorResponse(w, http.StatusNotFound, "no cached rates available")
		return
	}

	response := map[string]interface{}{
		"success": true,
		"rates":   rates,
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
