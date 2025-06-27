package handlers

import (
	"OtusGo/internal/application/services/exchange"
	"encoding/json"
	"log/slog"
	"net/http"
)

type ExchangeHandler struct {
	exchangeOrchestrator *exchange.ExchangeOrchestrator
	exchangeRateProvider *exchange.ExchangeRatesProvider
}

func NewExchangeHandler(exchangeOrchestrator *exchange.ExchangeOrchestrator,
	exchangeRateProvider *exchange.ExchangeRatesProvider) *ExchangeHandler {
	return &ExchangeHandler{
		exchangeOrchestrator: exchangeOrchestrator,
		exchangeRateProvider: exchangeRateProvider,
	}
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
	exchangeReq := exchange.NewExchangeOrchestratorRequest(
		req.UserID, req.FromCurrency, req.ToCurrency,
		req.FromAmount, req.CommissionRate)

	// Выполняем обмен
	result, err := h.exchangeOrchestrator.Exchange(r.Context(), exchangeReq)
	if err != nil {
		slog.Error("Exchange failed", "error", err)
		h.writeErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	// Формируем успешный ответ
	response := ExchangeResponseHTTP{
		Success:      true,
		ToAmount:     result.Amount,
		ExchangeRate: result.ExchangeRate,
		Commission:   result.Commission,
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

	rate, err := h.exchangeRateProvider.GetRate(r.Context(), fromCurrency, toCurrency)
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
