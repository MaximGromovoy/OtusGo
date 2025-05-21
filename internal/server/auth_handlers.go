package server

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// handleLogin обрабатывает запросы на аутентификацию
// @Summary Вход в систему
// @Description Аутентифицирует пользователя и выдает JWT токен
// @Tags авторизация
// @Accept json
// @Produce json
// @Param credentials body object{username=string,password=string} true "Учетные данные пользователя"
// @Success 200 {object} object{status=string,token=string} "Успешная аутентификация"
// @Failure 400 {string} string "Некорректный запрос"
// @Failure 401 {string} string "Неавторизован"
// @Router /login [post]
func (s *CurrencyServer) handleLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Метод не разрешен", http.StatusMethodNotAllowed)
		return
	}

	var credentials struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
		http.Error(w, "Некорректное тело запроса", http.StatusBadRequest)
		return
	}

	// В реальном приложении здесь была бы проверка учетных данных
	// Для демонстрации используем фиксированные значения
	if credentials.Username != "admin" || credentials.Password != "password" {
		http.Error(w, "Неверные учетные данные", http.StatusUnauthorized)
		return
	}

	// Создаем JWT токен
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": credentials.Username,
		"exp":      time.Now().Add(time.Hour * 24).Unix(), // Токен действителен 24 часа
	})

	// Подписываем токен нашим секретом
	tokenString, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		http.Error(w, "Не удалось сгенерировать токен", http.StatusInternalServerError)
		return
	}

	// Отправляем токен клиенту
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "успешно",
		"token":  tokenString,
	})
}
