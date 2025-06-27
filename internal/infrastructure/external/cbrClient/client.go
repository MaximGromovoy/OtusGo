package cbrClient

import (
	"OtusGo/internal/domain/models/currency"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"time"
)

// CBRService сервис для получения курсов валют с сайта ЦБ РФ
type CBRService struct {
	client *http.Client
}

// NewCBRClient создает новый экземпляр сервиса ЦБ РФ
func NewCBRClient() *CBRService {
	return &CBRService{
		client: &http.Client{
			Timeout: requestTimeout,
		},
	}
}

// GetCurrencyRates получает актуальные курсы валют с API ЦБ РФ
func (s *CBRService) GetCurrencyRates(ctx context.Context) ([]*currency.Currency, error) {
	cbrResponse, err := s.fetchCurrencyData(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch currency data: %w", err)
	}

	currencies := s.parseToCurrencyModels(cbrResponse)
	return currencies, nil
}

func (s *CBRService) GetSupportedCurrencies(ctx context.Context) ([]string, error) {
	cbrResponse, err := s.fetchCurrencyData(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch currency data: %w", err)
	}

	currencies := s.parseToCurrencyModels(cbrResponse)
	var codes []string
	for _, curr := range currencies {
		codes = append(codes, curr.GetCode())
	}

	return codes, nil
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
func (s *CBRService) parseToCurrencyModels(cbrResponse *CBRResponse) []*currency.Currency {
	var currencies []*currency.Currency

	for cbrCode, valuteData := range cbrResponse.Valute {
		// Нормализуем курс к единице валюты
		normalizedValue := valuteData.Value / float64(valuteData.Nominal)

		curr := currency.NewCurrency(valuteData.Name, valuteData.CharCode, normalizedValue)
		if curr != nil {
			currencies = append(currencies, curr)
			slog.Debug("Parsed currency",
				"code", cbrCode,
				"value", normalizedValue,
				"nominal", valuteData.Nominal,
				"original_value", valuteData.Value)
		} else {
			slog.Warn("Failed to create currency", "type", valuteData.CharCode)
		}
	}

	// Добавляем рубль как базовую валюту с курсом 1.0
	ruble := currency.NewCurrency("Рубль", "RUB", 1.0)
	if ruble != nil {
		currencies = append(currencies, ruble)
	}

	slog.Info("Parsed currencies", "count", len(currencies))
	return currencies
}
