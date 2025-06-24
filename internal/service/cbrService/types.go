package cbrService

import (
	"OtusGo/internal/interfaces"
	"net/http"
	"time"
)

type CBRResponse struct {
	Valute map[string]struct {
		CharCode string  `json:"CharCode"`
		Nominal  int     `json:"Nominal"`
		Value    float64 `json:"Value"`
		Name     string  `json:"Name"`
	} `json:"Valute"`
}

const (
	cbrAPIURL      = "https://www.cbr-xml-daily.ru/daily_json.js"
	requestTimeout = 10 * time.Second
	retryAttempts  = 3
	retryDelay     = 1 * time.Second
)

var _ interfaces.CBRServiceInterface = (*CBRService)(nil)

// CBRService сервис для получения курсов валют с сайта ЦБ РФ
type CBRService struct {
	client *http.Client
}

// NewCBRService создает новый экземпляр сервиса ЦБ РФ
func NewCBRService() *CBRService {
	return &CBRService{
		client: &http.Client{
			Timeout: requestTimeout,
		},
	}
}
