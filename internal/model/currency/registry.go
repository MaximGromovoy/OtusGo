package currency

type CurrencyConfig struct {
	Code string
}

var currencyConfigs = map[string]CurrencyConfig{
	DollarCurrency: {Code: DollarCode},
	EuroCurrency:   {Code: EuroCode},
	RubleCurrency:  {Code: RubleCode},
	LiraCurrency:   {Code: LiraCode},
}

func GetCurrencyConfig(currencyType string) (CurrencyConfig, bool) {
	config, exists := currencyConfigs[currencyType]
	return config, exists
}

func RegisterCurrencyConfig(currencyType, code string) {
	currencyConfigs[currencyType] = CurrencyConfig{
		Code: code,
	}

	RegisterCurrencyType(currencyType, 0)
}

func IsCurrencySupported(currencyType string) bool {
	_, exists := currencyConfigs[currencyType]
	return exists
}
