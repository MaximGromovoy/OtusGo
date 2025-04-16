package currency

type Currency struct {
	charCode string
	nominal  int
	value    float64
}

func NewCurrency(charCode string, nominal int, value float64) *Currency {
	return &Currency{
		charCode: charCode,
		nominal:  nominal,
		value:    value,
	}
}

func (c *Currency) GetCharCode() string {
	return c.charCode
}

func (c *Currency) GetValue() float64 {
	return c.value
}

func (c *Currency) GetNominal() int {
	return c.nominal
}
