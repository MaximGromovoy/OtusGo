package cbrService

import (
	"OtusGo/internal/model/currency"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"time"
)

// GetCurrencyRates получает актуальные курсы валют с API ЦБ РФ
func (s *CBRService) GetCurrencyRates(ctx context.Context) ([]currency.CurrencyInterface, error) {
	cbrResponse, err := s.fetchCurrencyData(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch currency data: %w", err)
	}

	currencies := s.parseToCurrencyModels(cbrResponse)
	return currencies, nil
}

// GetSupportedCurrencies возвращает список поддерживаемых валют
func (s *CBRService) GetSupportedCurrencies() []string {
	return []string{
		currency.DollarCurrency,
		currency.EuroCurrency,
		currency.LiraCurrency,
		currency.RubleCurrency,
	}
}

// fetchCurrencyData выполняет HTTP запрос к API ЦБ РФ с повторными попытками
func (s *CBRService) fetchCurrencyData(ctx context.Context) (*CBRResponse, error) {
	var lastErr error

	for attempt := 1; attempt <= retryAttempts; attempt++ {
		slog.Info("Fetching currency data from CBR", "attempt", attempt, "url", cbrAPIURL)

		req, err := http.NewRequestWithContext(ctx, http.MethodGet, cbrAPIURL, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to create request: %w", err)
		}

		req.Header.Set("Accept", "application/json")

		resp, err := s.client.Do(req)
		if err != nil {
			lastErr = fmt.Errorf("request failed: %w", err)
			slog.Warn("Request failed", "attempt", attempt, "error", err)

			if attempt < retryAttempts {
				select {
				case <-ctx.Done():
					return nil, ctx.Err()
				case <-time.After(retryDelay):
					continue
				}
			}
			continue
		}

		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			lastErr = fmt.Errorf("unexpected status code: %d", resp.StatusCode)
			slog.Warn("Unexpected status code", "attempt", attempt, "status", resp.StatusCode)

			if attempt < retryAttempts {
				select {
				case <-ctx.Done():
					return nil, ctx.Err()
				case <-time.After(retryDelay):
					continue
				}
			}
			continue
		}

		var cbrResponse CBRResponse
		if err := json.NewDecoder(resp.Body).Decode(&cbrResponse); err != nil {
			lastErr = fmt.Errorf("failed to decode response: %w", err)
			slog.Warn("Failed to decode response", "attempt", attempt, "error", err)

			if attempt < retryAttempts {
				select {
				case <-ctx.Done():
					return nil, ctx.Err()
				case <-time.After(retryDelay):
					continue
				}
			}
			continue
		}

		slog.Info("Successfully fetched currency data", "currencies_count", len(cbrResponse.Valute))
		return &cbrResponse, nil
	}

	return nil, fmt.Errorf("failed to fetch currency data after %d attempts: %w", retryAttempts, lastErr)
}

// parseToCurrencyModels конвертирует ответ ЦБ РФ в модели валют
func (s *CBRService) parseToCurrencyModels(cbrResponse *CBRResponse) []currency.CurrencyInterface {
	var currencies []currency.CurrencyInterface

	// Маппинг кодов валют ЦБ РФ на наши типы валют
	cbrToOurCurrency := map[string]string{
		currency.DollarCode: currency.DollarCurrency,
		currency.EuroCode:   currency.EuroCurrency,
		currency.LiraCode:   currency.LiraCurrency,
	}

	for cbrCode, valuteData := range cbrResponse.Valute {
		ourCurrencyType, exists := cbrToOurCurrency[cbrCode]
		if !exists {
			// Пропускаем неподдерживаемые валюты
			continue
		}

		// Нормализуем курс к единице валюты
		normalizedValue := valuteData.Value / float64(valuteData.Nominal)

		curr := currency.NewCurrency(ourCurrencyType, normalizedValue)
		if curr != nil {
			currencies = append(currencies, curr)
			slog.Debug("Parsed currency",
				"type", ourCurrencyType,
				"code", cbrCode,
				"value", normalizedValue,
				"nominal", valuteData.Nominal,
				"original_value", valuteData.Value)
		} else {
			slog.Warn("Failed to create currency", "type", ourCurrencyType)
		}
	}

	// Добавляем рубль как базовую валюту с курсом 1.0
	ruble := currency.NewCurrency(currency.RubleCurrency, 1.0)
	if ruble != nil {
		currencies = append(currencies, ruble)
	}

	slog.Info("Parsed currencies", "count", len(currencies))
	return currencies
}
