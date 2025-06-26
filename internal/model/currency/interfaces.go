package currency

type CurrencyInterface interface {
	GetName() string
	GetCode() string
	GetValue() float64
}
