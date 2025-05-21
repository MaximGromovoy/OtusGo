// filepath: d:\OTUS\GO\Rep\OtusGo\internal\server\utils.go
package server

import (
	"OtusGo/internal/model/currency"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

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
		return "", 0, fmt.Errorf("Not supported currency: %s", currencyType)
	}

	idStr := parts[3]
	if idStr == "" {
		return currencyType, 0, errors.New("Id must be specified")
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		return currencyType, 0, fmt.Errorf("Incorrect format: %w", err)
	}

	return currencyType, id, nil
}
