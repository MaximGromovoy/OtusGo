package currency

type BaseCurrency struct {
	name  string
	code  string
	value float64
	id    int
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

func (b *BaseCurrency) GetID() int {
	return b.id
}
