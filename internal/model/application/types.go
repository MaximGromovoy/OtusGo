package application

type Request struct {
	fromCurrency string
	toCurrency   string
	amount       float64
}

func NewRequest(fromCurrency, toCurrency string, amount float64) *Request {
	return &Request{
		fromCurrency: fromCurrency,
		toCurrency:   toCurrency,
		amount:       amount,
	}
}

func (r *Request) GetFromCurrency() string {
	return r.fromCurrency
}

func (r *Request) GetToCurrency() string {
	return r.toCurrency
}

func (r *Request) GetAmount() float64 {
	return r.amount
}

type Result struct {
	fromCurrency    string
	toCurrency      string
	originalAmount  float64
	convertedAmount float64
	rate            float64
}

func NewResult(fromCurrency, toCurrency string, originalAmount, convertedAmount, rate float64) *Result {
	return &Result{
		fromCurrency:    fromCurrency,
		toCurrency:      toCurrency,
		originalAmount:  originalAmount,
		convertedAmount: convertedAmount,
		rate:            rate,
	}
}

func (r *Result) GetFromCurrency() string {
	return r.fromCurrency
}

func (r *Result) GetToCurrency() string {
	return r.toCurrency
}

func (r *Result) GetOriginalAmount() float64 {
	return r.originalAmount
}

func (r *Result) GetConvertedAmount() float64 {
	return r.convertedAmount
}

func (r *Result) GetRate() float64 {
	return r.rate
}
