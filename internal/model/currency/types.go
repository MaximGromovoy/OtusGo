package currency

type CurrencyInterface interface {
	GetName() string
	GetCode() string
	GetValue() float64
}

type BaseCurrency struct {
	name  string
	code  string
	value float64
}

func (b *BaseCurrency) GetName() string {
	return b.name
}

func (b *BaseCurrency) GetCode() string {
	return b.code
}

func (b *BaseCurrency) GetValue() float64 {
	return b.value
}

type Ruble struct {
	BaseCurrency
}

// Конструктор для рубля
func NewRuble(value float64) *Ruble {
	return &Ruble{
		BaseCurrency: BaseCurrency{
			name:  "Ruble",
			code:  "RUB",
			value: value,
		},
	}
}

type Dollar struct {
	BaseCurrency
}

func NewDollar(value float64) *Dollar {
	return &Dollar{
		BaseCurrency: BaseCurrency{
			name:  "Dollar",
			code:  "USD",
			value: value,
		},
	}
}

type Euro struct {
	BaseCurrency
}

func NewEuro(value float64) *Euro {
	return &Euro{
		BaseCurrency: BaseCurrency{
			name:  "Euro",
			code:  "EUR",
			value: value,
		},
	}
}

type Lira struct {
	BaseCurrency
}

func NewLira(value float64) *Lira {
	return &Lira{
		BaseCurrency: BaseCurrency{
			name:  "Lira",
			code:  "TRY",
			value: value,
		},
	}
}
