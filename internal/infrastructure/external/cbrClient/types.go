package cbrClient

type CBRResponse struct {
	Valute map[string]struct {
		CharCode string  `json:"CharCode"`
		Nominal  int     `json:"Nominal"`
		Value    float64 `json:"Value"`
		Name     string  `json:"Name"`
	} `json:"Valute"`
}
