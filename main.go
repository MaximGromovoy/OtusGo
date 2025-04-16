package main

import (
	"OtusGo/internal/model/currency"
	"OtusGo/internal/repository"
	"OtusGo/internal/service"
	"fmt"
)

func main() {
	var currencies []currency.CurrencyInterface = service.GetRandomCurrencies()

	rubles, dollars, euros, liras := repository.DistributeCurrencies(currencies)

	fmt.Println("Rubles:")
	for _, ruble := range rubles {
		fmt.Printf("Name: %s, Code: %s, Value: %.2f\n", ruble.GetName(), ruble.GetCode(), ruble.GetValue())
	}

	fmt.Println("\nDollars:")
	for _, dollar := range dollars {
		fmt.Printf("Name: %s, Code: %s, Value: %.2f\n", dollar.GetName(), dollar.GetCode(), dollar.GetValue())
	}

	fmt.Println("\nEuros:")
	for _, euro := range euros {
		fmt.Printf("Name: %s, Code: %s, Value: %.2f\n", euro.GetName(), euro.GetCode(), euro.GetValue())
	}

	fmt.Println("\nLiras:")
	for _, lira := range liras {
		fmt.Printf("Name: %s, Code: %s, Value: %.2f\n", lira.GetName(), lira.GetCode(), lira.GetValue())
	}
}
