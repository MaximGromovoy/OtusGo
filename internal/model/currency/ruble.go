package currency

type Ruble struct {
	BaseCurrency
}

func NewRuble(value float64) *Ruble {
	return &Ruble{
		BaseCurrency: BaseCurrency{
			name:  RubleCurrency,
			code:  "RUB",
			value: value,
		},
	}
}
