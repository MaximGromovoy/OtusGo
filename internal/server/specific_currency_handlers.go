package server

import (
	"OtusGo/internal/model/currency"
	"encoding/json"
	"fmt"
	"net/http"
)

// handleCurrencyWithTypeAndID обрабатывает GET, PUT, DELETE запросы для конкретной валюты по типу и ID
func (s *CurrencyServer) handleCurrencyWithTypeAndID(w http.ResponseWriter, r *http.Request) {
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
// @Summary Получение валюты
// @Description Получает информацию о конкретной валюте по типу и ID
// @Tags валюта
// @Produce json
// @Param type path string true "Тип валюты" Enums(Dollar, Euro, Ruble, Lira)
// @Param id path integer true "ID валюты"
// @Success 200 {object} object{status=string,data=object{id=integer,type=string,code=string,value=number}} "Информация о валюте"
// @Failure 400 {string} string "Некорректный запрос"
// @Failure 401 {string} string "Не авторизован"
// @Failure 404 {string} string "Не найдено"
// @Security JWT
// @Router /currency/{type}/{id} [get]
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
// @Summary Обновление валюты
// @Description Обновляет значение валюты указанного типа по ID
// @Tags валюта
// @Accept json
// @Produce json
// @Param type path string true "Тип валюты" Enums(Dollar, Euro, Ruble, Lira)
// @Param id path integer true "ID валюты"
// @Param currency body object{value=number} true "Новое значение валюты"
// @Success 200 {object} object{status=string,data=object{id=integer,type=string,code=string,value=number}} "Валюта успешно обновлена"
// @Failure 400 {string} string "Некорректный запрос"
// @Failure 401 {string} string "Не авторизован"
// @Failure 404 {string} string "Не найдено"
// @Failure 500 {string} string "Внутренняя ошибка сервера"
// @Security JWT
// @Router /currency/{type}/{id} [put]
func (s *CurrencyServer) updateCurrency(w http.ResponseWriter, r *http.Request, currencyType string, id int) {
	// Получаем значение валюты из тела запроса
	var requestData struct {
		Value float64 `json:"value"`
	}

	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
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
// @Summary Удаление валюты
// @Description Удаляет валюту указанного типа по ID
// @Tags валюта
// @Produce json
// @Param type path string true "Тип валюты" Enums(Dollar, Euro, Ruble, Lira)
// @Param id path integer true "ID валюты"
// @Success 200 {object} object{status=string,message=string} "Валюта успешно удалена"
// @Failure 400 {string} string "Некорректный запрос"
// @Failure 401 {string} string "Не авторизован"
// @Failure 500 {string} string "Внутренняя ошибка сервера"
// @Security JWT
// @Router /currency/{type}/{id} [delete]
func (s *CurrencyServer) deleteCurrency(w http.ResponseWriter, currencyType string, id int) {
	// Удаляем из репозитория
	if err := s.repo.DeleteByTypeAndID(currencyType, id); err != nil {
		http.Error(w, fmt.Sprintf("Failed to delete currency: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  "success",
		"message": fmt.Sprintf("Currency with type %s and ID %d successfully deleted", currencyType, id),
	})
}
