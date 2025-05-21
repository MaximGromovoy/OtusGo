package server

import (
	"OtusGo/internal/model/currency"
	"OtusGo/internal/repository/currencyRepository"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type CurrencyServer struct {
	repo *currencyRepository.CurrencyRepository
}

func NewCurrencyServer(repo *currencyRepository.CurrencyRepository) *CurrencyServer {
	return &CurrencyServer{
		repo: repo,
	}
}

// Start инициализирует и запускает HTTP сервер
func (s *CurrencyServer) Start(port string) error {
	http.HandleFunc("/currency", s.handleCurrency)
	http.HandleFunc("/currencies", s.handleCurrencies)
	http.HandleFunc("/currency/", s.handleCurrencyWithTypeAndID)

	log.Printf("Starting server on port %s", port)
	return http.ListenAndServe(":"+port, nil)
}

// parseRequestedCurrencyTypeAndID извлекает тип валюты и ID из URL вида /currency/{type}/{id}
func (s *CurrencyServer) parseRequestedCurrencyTypeAndID(path string) (string, int, error) {
	parts := strings.Split(path, "/")
	if len(parts) != 4 {
		return "", 0, errors.New("invalid URL format, expected /currency/{type}/{id}")
	}

	currencyType := parts[2]
	if currencyType == "" {
		return "", 0, errors.New("currency type is required")
	}

	if !currency.IsCurrencySupported(currencyType) {
		return "", 0, fmt.Errorf("unsupported currency type: %s", currencyType)
	}

	idStr := parts[3]
	if idStr == "" {
		return currencyType, 0, errors.New("currency ID is required")
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		return currencyType, 0, fmt.Errorf("invalid ID format: %w", err)
	}

	return currencyType, id, nil
}

// handleCurrency обрабатывает POST запросы для добавления новой валюты
func (s *CurrencyServer) handleCurrency(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		s.addCurrency(w, r)
		return
	}

	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
}

// addCurrency добавляет новую валюту
func (s *CurrencyServer) addCurrency(w http.ResponseWriter, r *http.Request) {
	// Получаем тип валюты из параметра запроса
	currencyType := r.URL.Query().Get("type")
	if currencyType == "" {
		http.Error(w, "Currency type is required as a query parameter", http.StatusBadRequest)
		return
	}

	// Проверяем поддерживается ли тип валюты
	if !currency.IsCurrencySupported(currencyType) {
		http.Error(w, fmt.Sprintf("Unsupported currency type: %s", currencyType), http.StatusBadRequest)
		return
	}

	// Получаем значение валюты из тела запроса
	var requestData struct {
		Value float64 `json:"value"`
	}

	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		http.Error(w, "Invalid request body, expected {\"value\": 123.45}", http.StatusBadRequest)
		return
	}

	// Создаем новую валюту
	newCurrency := currency.NewCurrency(currencyType, requestData.Value)
	if newCurrency == nil {
		http.Error(w, "Failed to create currency", http.StatusInternalServerError)
		return
	}

	// Добавляем в репозиторий
	if err := s.repo.Add(newCurrency); err != nil {
		http.Error(w, fmt.Sprintf("Failed to add currency: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "success",
		"data": map[string]interface{}{
			"id":    newCurrency.GetID(),
			"type":  newCurrency.GetName(),
			"code":  newCurrency.GetCode(),
			"value": newCurrency.GetValue(),
		},
	})
}

// handleCurrencies обрабатывает GET запросы для получения всех валют
func (s *CurrencyServer) handleCurrencies(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Получаем тип валюты из параметра запроса (необязательный)
	requestedType := r.URL.Query().Get("type")

	result := make(map[string][]map[string]interface{})

	// Если тип указан и поддерживается, возвращаем только его
	if requestedType != "" {
		if !currency.IsCurrencySupported(requestedType) {
			http.Error(w, fmt.Sprintf("Unsupported currency type: %s", requestedType), http.StatusBadRequest)
			return
		}

		currencies := s.repo.GetAll(requestedType)
		currenciesData := make([]map[string]interface{}, 0, len(currencies))

		for _, curr := range currencies {
			currenciesData = append(currenciesData, map[string]interface{}{
				"id":    curr.GetID(),
				"type":  curr.GetName(),
				"code":  curr.GetCode(),
				"value": curr.GetValue(),
			})
		}

		result[requestedType] = currenciesData
	} else {
		// Иначе возвращаем все поддерживаемые типы
		for _, currType := range currency.SupportedCurrencies {
			currencies := s.repo.GetAll(currType)
			currenciesData := make([]map[string]interface{}, 0, len(currencies))

			for _, curr := range currencies {
				currenciesData = append(currenciesData, map[string]interface{}{
					"id":    curr.GetID(),
					"type":  curr.GetName(),
					"code":  curr.GetCode(),
					"value": curr.GetValue(),
				})
			}

			result[currType] = currenciesData
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "success",
		"data":   result,
	})
}

// handleCurrencyWithTypeAndID обрабатывает GET, PUT, DELETE запросы для конкретной валюты по типу и ID
func (s *CurrencyServer) handleCurrencyWithTypeAndID(w http.ResponseWriter, r *http.Request) {
	// Извлекаем тип и ID из URL
	currencyType, id, err := s.parseRequestedCurrencyTypeAndID(r.URL.Path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		s.getCurrency(w, currencyType, id)
	case http.MethodPut:
		s.updateCurrency(w, r, currencyType, id)
	case http.MethodDelete:
		s.deleteCurrency(w, currencyType, id)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// getCurrency получает валюту по типу и ID
func (s *CurrencyServer) getCurrency(w http.ResponseWriter, currencyType string, id int) {
	currencies := s.repo.GetAll(currencyType)
	for _, curr := range currencies {
		if curr.GetID() == id {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]interface{}{
				"status": "success",
				"data": map[string]interface{}{
					"id":    curr.GetID(),
					"type":  curr.GetName(),
					"code":  curr.GetCode(),
					"value": curr.GetValue(),
				},
			})
			return
		}
	}

	http.Error(w, fmt.Sprintf("Currency with type %s and ID %d not found", currencyType, id), http.StatusNotFound)
}

// updateCurrency обновляет валюту по типу и ID
func (s *CurrencyServer) updateCurrency(w http.ResponseWriter, r *http.Request, currencyType string, id int) {
	// Получаем значение валюты из тела запроса
	var requestData struct {
		Value float64 `json:"value"`
	}

	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		http.Error(w, "Invalid request body, expected {\"value\": 123.45}", http.StatusBadRequest)
		return
	}

	// Находим валюту для обновления
	currencies := s.repo.GetAll(currencyType)
	var currToUpdate currency.CurrencyInterface

	for _, curr := range currencies {
		if curr.GetID() == id {
			currToUpdate = curr
			break
		}
	}

	if currToUpdate == nil {
		http.Error(w, fmt.Sprintf("Currency with type %s and ID %d not found", currencyType, id), http.StatusNotFound)
		return
	}

	// Создаем обновленную валюту
	updatedCurrency := currency.NewCurrencyWithID(currencyType, requestData.Value, id)
	if updatedCurrency == nil {
		http.Error(w, "Failed to create updated currency", http.StatusInternalServerError)
		return
	}

	// Обновляем в репозитории
	if err := s.repo.Update(currencyType, id, updatedCurrency); err != nil {
		http.Error(w, fmt.Sprintf("Failed to update currency: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "success",
		"data": map[string]interface{}{
			"id":    updatedCurrency.GetID(),
			"type":  updatedCurrency.GetName(),
			"code":  updatedCurrency.GetCode(),
			"value": updatedCurrency.GetValue(),
		},
	})
}

// deleteCurrency удаляет валюту по типу и ID
func (s *CurrencyServer) deleteCurrency(w http.ResponseWriter, currencyType string, id int) {
	// Удаляем из репозитория
	if err := s.repo.DeleteByTypeAndID(currencyType, id); err != nil {
		http.Error(w, fmt.Sprintf("Failed to delete currency: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  "success",
		"message": fmt.Sprintf("Currency with type %s and ID %d deleted successfully", currencyType, id),
	})
}
