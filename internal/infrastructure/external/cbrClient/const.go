package cbrClient

import "time"

const (
	cbrAPIURL      = "https://www.cbr-xml-daily.ru/daily_json.js"
	requestTimeout = 10 * time.Second
	retryAttempts  = 3
	retryDelay     = 1 * time.Second
)
