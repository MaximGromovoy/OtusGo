package currency

func NewCurrency(currencyType string, value float64) CurrencyInterface {
	config, exists := GetCurrencyConfig(currencyType)
	if !exists {
		return nil
	}

	base := BaseCurrency{
		name:  currencyType,
		code:  config.Code,
		value: value,
	}

	return createCurrency(currencyType, base)
}

func NewCurrencyWithID(currencyType string, value float64, id int) CurrencyInterface {
	config, exists := GetCurrencyConfig(currencyType)
	if !exists {
		return nil
	}

	base := BaseCurrency{
		name:  currencyType,
		code:  config.Code,
		value: value,
	}

	return createCurrency(currencyType, base)
}

func createCurrency(currencyType string, base BaseCurrency) CurrencyInterface {
	switch currencyType {
	case DollarCurrency:
		return &Dollar{BaseCurrency: base}
	case EuroCurrency:
		return &Euro{BaseCurrency: base}
	case RubleCurrency:
		return &Ruble{BaseCurrency: base}
	case LiraCurrency:
		return &Lira{BaseCurrency: base}
	default:
		return nil
	}
}
