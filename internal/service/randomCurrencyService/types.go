package randomCurrencyService

import (
	"OtusGo/internal/model/currency"
	"math/rand"
	"sync"
	"time"
)

type RandomCurrencyService struct {
	randomSeed rand.Rand
}

func NewRandomCurrencyService() *RandomCurrencyService {
	return &RandomCurrencyService{
		randomSeed: *rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

func (randomCurrencyService *RandomCurrencyService) GetRandomCurrencies() (currencies []currency.CurrencyInterface) {
	for range 20 {
		currencies = append(currencies, randomCurrencyService.getRandomCurrency())
	}
	return
}

func (randomCurrencyService *RandomCurrencyService) StartGenerateCurrencies(currencyChannel chan currency.CurrencyInterface, wg *sync.WaitGroup) {
	wg.Add(1)

	go func() {
		defer wg.Done()

		ticker := time.NewTicker(50 * time.Millisecond)
		defer ticker.Stop()

		timeout := time.After(1 * time.Second)

		for {
			select {
			case <-timeout:
				close(currencyChannel)
				return
			case <-ticker.C:
				for range 1 {
					currencyChannel <- randomCurrencyService.getRandomCurrency()
				}
			}
		}
	}()
}

func (randomCurrencyService *RandomCurrencyService) getRandomCurrency() currency.CurrencyInterface {
	r := randomCurrencyService.randomSeed

	localCurrencies := []currency.CurrencyInterface{
		currency.NewRuble(r.Float64() * 100),
		currency.NewDollar(r.Float64() * 100),
		currency.NewEuro(r.Float64() * 100),
		currency.NewLira(r.Float64() * 100),
	}

	return localCurrencies[r.Intn(len(localCurrencies))]
}
