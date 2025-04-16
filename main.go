package main

import (
	"OtusGo/internal/repository/currencyRepository"
	"OtusGo/internal/service"
	"fmt"
)

func main() {
	repository := currencyRepository.NewCurrencyRepository()
	var currencies = service.GetRandomCurrencies()
	repository.AddCurrencies(currencies)
	LogCurrencies(repository)
}

func LogCurrencies(repository *currencyRepository.CurrencyRepository) {
	fmt.Println("Rubles:")
	for _, ruble := range repository.GetRubles() {
		fmt.Printf("Name: %s, Code: %s, Value: %.2f\n", ruble.GetName(), ruble.GetCode(), ruble.GetValue())
	}

	fmt.Println("\nDollars:")
	for _, dollar := range repository.GetDollars() {
		fmt.Printf("Name: %s, Code: %s, Value: %.2f\n", dollar.GetName(), dollar.GetCode(), dollar.GetValue())
	}

	fmt.Println("\nEuros:")
	for _, euro := range repository.GetEuros() {
		fmt.Printf("Name: %s, Code: %s, Value: %.2f\n", euro.GetName(), euro.GetCode(), euro.GetValue())
	}

	fmt.Println("\nLiras:")
	for _, lira := range repository.GetLiras() {
		fmt.Printf("Name: %s, Code: %s, Value: %.2f\n", lira.GetName(), lira.GetCode(), lira.GetValue())
	}
}
