package postgresRepository

import (
	"fmt"
)

type Configuration struct {
	uri string
}

func NewConfigurationFromParts(address, dbName, sslMode string) *Configuration {
	return &Configuration{
		uri: createUri(address, dbName, sslMode),
	}
}

func NewConfigurationFromUri(uri string) *Configuration {
	return &Configuration{
		uri: uri,
	}
}

func createUri(address, dbName, sslMode string) string {
	return fmt.Sprintf("%s/%s?%s", address, dbName, sslMode)
}
