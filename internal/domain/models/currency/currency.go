package currency

type Currency struct {
	name  string
	code  string
	value float64
}

func NewCurrency(name, code string, value float64) *Currency {
	return &Currency{
		name:  name,
		code:  code,
		value: value,
	}
}

func (b *Currency) GetName() string {
	return b.name
}

func (b *Currency) GetCode() string {
	return b.code
}

func (b *Currency) GetValue() float64 {
	return b.value
}
