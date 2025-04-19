package currency

type Euro struct {
	BaseCurrency
}

func NewEuro(value float64) *Euro {
	return &Euro{
		BaseCurrency: BaseCurrency{
			name:  EuroCurrency,
			code:  "EUR",
			value: value,
		},
	}
}
