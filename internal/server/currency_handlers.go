package server

import (
	"OtusGo/internal/model/currency"
	"encoding/json"
	"fmt"
	"net/http"
)

// handleCurrency обрабатывает POST запросы для добавления новой валюты
// @Summary Добавление валюты
// @Description Добавляет новую валюту указанного типа
// @Tags валюта
// @Accept json
// @Produce json
// @Param type query string true "Тип валюты" Enums(Dollar, Euro, Ruble, Lira)
// @Param currency body object{value=number} true "Данные валюты"
// @Success 201 {object} object{status=string,data=object{id=integer,type=string,code=string,value=number}} "Валюта успешно добавлена"
// @Failure 400 {string} string "Некорректный запрос"
// @Failure 401 {string} string "Не авторизован"
// @Failure 500 {string} string "Внутренняя ошибка сервера"
// @Security JWT
// @Router /currency [post]
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
		http.Error(w, "Currency type is required", http.StatusBadRequest)
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
		http.Error(w, "Invalid request body", http.StatusBadRequest)
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
// @Summary Получение списка валют
// @Description Возвращает список всех валют или валют определенного типа
// @Tags валюта
// @Produce json
// @Param type query string false "Тип валюты (опционально)" Enums(Dollar, Euro, Ruble, Lira)
// @Success 200 {object} object{status=string,data=object} "Список валют"
// @Failure 400 {string} string "Некорректный запрос"
// @Failure 401 {string} string "Не авторизован"
// @Security JWT
// @Router /currencies [get]
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
