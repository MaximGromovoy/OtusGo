package exchange

import (
	"OtusGo/internal/domain/services/exchangeService"
	"context"
	"strings"
	"testing"
)

func setupOrchestratorForValidation() *ExchangeOrchestrator {
	provider := &mockExchangeRatesProvider{
		rates: map[string]float64{
			"USD-EUR": 0.85,
		},
	}
	repo := newMockTransactionRepository()
	service := exchangeService.NewExchangeService()
	logger := &mockLogger{}

	return NewExchangeOrchestrator(provider, repo, service, logger)

}

func TestExchangeOrchestrator_Exchange_NilRequest(t *testing.T) {
	// Arrange
	orchestrator := setupOrchestratorForValidation()
	ctx := context.Background()

	// Act
	result, err := orchestrator.Exchange(ctx, nil)

	// Assert
	if err == nil {
		t.Fatal("Expected validation error for nil request")
	}
	if result != nil {
		t.Error("Expected nil result on validation error")
	}
	expectedError := "validation failed: request cannot be nil"
	if err.Error() != expectedError {
		t.Errorf("Expected error '%s', got '%s'", expectedError, err.Error())
	}
}

func TestExchangeOrchestrator_Exchange_InvalidUserID(t *testing.T) {
	testCases := []struct {
		name     string
		userID   int
		expected string
	}{
		{
			name:     "Zero user ID",
			userID:   0,
			expected: "validation failed: user ID must be positive, got 0",
		},
		{
			name:     "Negative user ID",
			userID:   -1,
			expected: "validation failed: user ID must be positive, got -1",
		},
	}

	orchestrator := setupOrchestratorForValidation()
	ctx := context.Background()

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			req := NewExchangeOrchestratorRequest(tc.userID, "USD", "EUR", 100.0, 2.0)

			// Act
			result, err := orchestrator.Exchange(ctx, req)

			// Assert
			if err == nil {
				t.Fatal("Expected validation error for invalid user ID")
			}
			if result != nil {
				t.Error("Expected nil result on validation error")
			}
			if err.Error() != tc.expected {
				t.Errorf("Expected error '%s', got '%s'", tc.expected, err.Error())
			}
		})
	}
}

func TestExchangeOrchestrator_Exchange_EmptyCurrency(t *testing.T) {
	testCases := []struct {
		name         string
		fromCurrency string
		toCurrency   string
		expected     string
	}{
		{
			name:         "Empty from currency",
			fromCurrency: "",
			toCurrency:   "EUR",
			expected:     "validation failed: from currency cannot be empty",
		},
		{
			name:         "Empty to currency",
			fromCurrency: "USD",
			toCurrency:   "",
			expected:     "validation failed: to currency cannot be empty",
		},
		{
			name:         "Whitespace from currency",
			fromCurrency: "   ",
			toCurrency:   "EUR",
			expected:     "validation failed: from currency cannot be empty",
		},
		{
			name:         "Whitespace to currency",
			fromCurrency: "USD",
			toCurrency:   "   ",
			expected:     "validation failed: to currency cannot be empty",
		},
	}

	orchestrator := setupOrchestratorForValidation()
	ctx := context.Background()

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			req := NewExchangeOrchestratorRequest(1, tc.fromCurrency, tc.toCurrency, 100.0, 2.0)

			// Act
			result, err := orchestrator.Exchange(ctx, req)

			// Assert
			if err == nil {
				t.Fatal("Expected validation error for empty currency")
			}
			if result != nil {
				t.Error("Expected nil result on validation error")
			}
			if err.Error() != tc.expected {
				t.Errorf("Expected error '%s', got '%s'", tc.expected, err.Error())
			}
		})
	}
}

func TestExchangeOrchestrator_Exchange_SameCurrencies(t *testing.T) {
	// Arrange
	orchestrator := setupOrchestratorForValidation()
	ctx := context.Background()
	req := NewExchangeOrchestratorRequest(1, "USD", "USD", 100.0, 2.0)

	// Act
	result, err := orchestrator.Exchange(ctx, req)

	// Assert
	if err == nil {
		t.Fatal("Expected validation error for same currencies")
	}
	if result != nil {
		t.Error("Expected nil result on validation error")
	}
	expectedError := "validation failed: from currency and to currency cannot be the same: USD"
	if err.Error() != expectedError {
		t.Errorf("Expected error '%s', got '%s'", expectedError, err.Error())
	}
}

func TestExchangeOrchestrator_Exchange_NegativeAmount(t *testing.T) {
	// Arrange
	orchestrator := setupOrchestratorForValidation()
	ctx := context.Background()
	req := NewExchangeOrchestratorRequest(1, "USD", "EUR", -100.0, 2.0)

	// Act
	result, err := orchestrator.Exchange(ctx, req)

	// Assert
	if err == nil {
		t.Fatal("Expected validation error for negative amount")
	}
	if result != nil {
		t.Error("Expected nil result on validation error")
	}
	expectedError := "validation failed: amount cannot be negative, got -100.000000"
	if err.Error() != expectedError {
		t.Errorf("Expected error '%s', got '%s'", expectedError, err.Error())
	}
}

func TestExchangeOrchestrator_Exchange_InvalidCommissionRate(t *testing.T) {
	testCases := []struct {
		name           string
		commissionRate float64
		expected       string
	}{
		{
			name:           "Negative commission rate",
			commissionRate: -1.0,
			expected:       "validation failed: commission rate cannot be negative, got -1.000000",
		},
		{
			name:           "Commission rate 100%",
			commissionRate: 100.0,
			expected:       "validation failed: commission rate cannot be 100% or more, got 100.000000",
		},
		{
			name:           "Commission rate over 100%",
			commissionRate: 150.0,
			expected:       "validation failed: commission rate cannot be 100% or more, got 150.000000",
		},
	}

	orchestrator := setupOrchestratorForValidation()
	ctx := context.Background()

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			req := NewExchangeOrchestratorRequest(1, "USD", "EUR", 100.0, tc.commissionRate)

			// Act
			result, err := orchestrator.Exchange(ctx, req)

			// Assert
			if err == nil {
				t.Fatal("Expected validation error for invalid commission rate")
			}
			if result != nil {
				t.Error("Expected nil result on validation error")
			}
			if err.Error() != tc.expected {
				t.Errorf("Expected error '%s', got '%s'", tc.expected, err.Error())
			}
		})
	}
}

func TestExchangeOrchestrator_Exchange_InvalidCurrencyCodeLength(t *testing.T) {
	testCases := []struct {
		name         string
		fromCurrency string
		toCurrency   string
		expected     string
	}{
		{
			name:         "From currency too short",
			fromCurrency: "US",
			toCurrency:   "EUR",
			expected:     "validation failed: from currency must be 3 characters long, got US",
		},
		{
			name:         "From currency too long",
			fromCurrency: "USDX",
			toCurrency:   "EUR",
			expected:     "validation failed: from currency must be 3 characters long, got USDX",
		},
		{
			name:         "To currency too short",
			fromCurrency: "USD",
			toCurrency:   "EU",
			expected:     "validation failed: to currency must be 3 characters long, got EU",
		},
		{
			name:         "To currency too long",
			fromCurrency: "USD",
			toCurrency:   "EURO",
			expected:     "validation failed: to currency must be 3 characters long, got EURO",
		},
	}

	orchestrator := setupOrchestratorForValidation()
	ctx := context.Background()

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			req := NewExchangeOrchestratorRequest(1, tc.fromCurrency, tc.toCurrency, 100.0, 2.0)

			// Act
			result, err := orchestrator.Exchange(ctx, req)

			// Assert
			if err == nil {
				t.Fatal("Expected validation error for invalid currency code length")
			}
			if result != nil {
				t.Error("Expected nil result on validation error")
			}
			if err.Error() != tc.expected {
				t.Errorf("Expected error '%s', got '%s'", tc.expected, err.Error())
			}
		})
	}
}

func TestExchangeOrchestrator_Exchange_ValidZeroAmount(t *testing.T) {
	// Arrange
	orchestrator := setupOrchestratorForValidation()
	ctx := context.Background()
	req := NewExchangeOrchestratorRequest(1, "USD", "EUR", 0.0, 2.0)

	// Act
	result, err := orchestrator.Exchange(ctx, req)

	// Assert
	if err != nil {
		t.Fatalf("Expected no error for zero amount, got %v", err)
	}
	if result == nil {
		t.Fatal("Expected result for valid zero amount")
	}
	if result.Amount != 0.0 {
		t.Errorf("Expected amount 0.0, got %f", result.Amount)
	}
}

func TestExchangeOrchestrator_Exchange_ValidMaxCommissionRate(t *testing.T) {
	// Arrange
	orchestrator := setupOrchestratorForValidation()
	ctx := context.Background()
	req := NewExchangeOrchestratorRequest(1, "USD", "EUR", 100.0, 99.99)

	// Act
	result, err := orchestrator.Exchange(ctx, req)

	// Assert
	if err != nil {
		t.Fatalf("Expected no error for max valid commission rate, got %v", err)
	}
	if result == nil {
		t.Fatal("Expected result for valid commission rate")
	}
}

func TestExchangeOrchestrator_Exchange_CurrencyCodeWithSpaces(t *testing.T) {
	// Arrange
	orchestrator := setupOrchestratorForValidation()
	ctx := context.Background()

	// Валидный запрос с пробелами в валютных кодах (должны быть обрезаны)
	req := NewExchangeOrchestratorRequest(1, " USD ", " EUR ", 100.0, 2.0)

	// Act
	result, err := orchestrator.Exchange(ctx, req)

	// Assert
	if err != nil {
		t.Fatalf("Expected no error for currency codes with spaces, got %v", err)
	}
	if result == nil {
		t.Fatal("Expected result for valid currency codes with spaces")
	}
}

func TestExchangeOrchestrator_validateRequest_DirectCall(t *testing.T) {
	// Arrange
	orchestrator := setupOrchestratorForValidation()

	testCases := []struct {
		name        string
		req         *ExchangeOrchestratorRequest
		expectError bool
		errorMsg    string
	}{
		{
			name:        "Valid request",
			req:         NewExchangeOrchestratorRequest(1, "USD", "EUR", 100.0, 2.0),
			expectError: false,
		},
		{
			name:        "Nil request",
			req:         nil,
			expectError: true,
			errorMsg:    "request cannot be nil",
		},
		{
			name:        "Multiple validation errors - should return first",
			req:         NewExchangeOrchestratorRequest(-1, "", "", -100.0, -5.0),
			expectError: true,
			errorMsg:    "user ID must be positive, got -1",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Act
			err := orchestrator.validateRequest(tc.req)

			// Assert
			if tc.expectError && err == nil {
				t.Fatal("Expected validation error, got nil")
			}
			if !tc.expectError && err != nil {
				t.Fatalf("Expected no validation error, got %v", err)
			}
			if tc.expectError && tc.errorMsg != "" && !strings.Contains(err.Error(), tc.errorMsg) {
				t.Errorf("Expected error to contain '%s', got '%s'", tc.errorMsg, err.Error())
			}
		})
	}
}
