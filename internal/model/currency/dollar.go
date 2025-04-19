package currency

type Dollar struct {
	BaseCurrency
}

func NewDollar(value float64) *Dollar {
	return &Dollar{
		BaseCurrency: BaseCurrency{
			name:  DollarCurrency,
			code:  "USD",
			value: value,
		},
	}
}
