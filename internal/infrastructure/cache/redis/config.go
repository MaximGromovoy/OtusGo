package redisCache

type Configuration struct {
	address  string
	password string
	dbNumber int
}

func NewConfiguration(address, password string, dbNumber int) *Configuration {
	return &Configuration{
		address:  address,
		password: password,
		dbNumber: dbNumber,
	}
}

func (config *Configuration) SetAddress(address string) {
	config.address = address
}

func (config *Configuration) SetPassword(password string) {
	config.password = password
}

func (config *Configuration) SetDbNumber(dbNumber int) {
	config.dbNumber = dbNumber
}
