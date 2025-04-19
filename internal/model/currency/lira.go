package currency

type Lira struct {
	BaseCurrency
}

func NewLira(value float64) *Lira {
	return &Lira{
		BaseCurrency: BaseCurrency{
			name:  LiraCurrency,
			code:  "TRY",
			value: value,
		},
	}
}
