package currencyRepositoryWatcher

import (
	"OtusGo/internal/model/currency"
	"OtusGo/internal/repository/currencyRepository"
	"fmt"
	"sync"
	"time"
)

type CurrencyRepositoryWatcher struct {
	currencyRepository *currencyRepository.CurrencyRepository
	prevRubles         []*currency.Ruble
	prevDollars        []*currency.Dollar
	prevEuros          []*currency.Euro
	prevLiras          []*currency.Lira
	stopChannel        chan struct{}
	wg                 sync.WaitGroup
}

func NewCurrencyRepositoryWatcher(currencyRepository *currencyRepository.CurrencyRepository) *CurrencyRepositoryWatcher {
	return &CurrencyRepositoryWatcher{
		currencyRepository: currencyRepository,
		prevRubles:         make([]*currency.Ruble, 0),
		prevDollars:        make([]*currency.Dollar, 0),
		prevEuros:          make([]*currency.Euro, 0),
		prevLiras:          make([]*currency.Lira, 0),
		stopChannel:        make(chan struct{}),
	}
}

func (watcher *CurrencyRepositoryWatcher) StartWatching() {
	watcher.stopChannel = make(chan struct{})
	watcher.wg.Add(1)
	go func() {
		defer watcher.wg.Done()

		ticker := time.NewTicker(200 * time.Millisecond)
		defer ticker.Stop()

		for {
			select {
			case <-watcher.stopChannel:
				fmt.Println("Stop channel was closed. Goroutine was stopped.")
				return
			case <-ticker.C:
				currentRubles := watcher.currencyRepository.GetRubles()
				currentDollars := watcher.currencyRepository.GetDollars()
				currentEuros := watcher.currencyRepository.GetEuros()
				currentLiras := watcher.currencyRepository.GetLiras()

				logAddedCurrencies(currentRubles, watcher.prevRubles)
				logAddedCurrencies(currentDollars, watcher.prevDollars)
				logAddedCurrencies(currentEuros, watcher.prevEuros)
				logAddedCurrencies(currentLiras, watcher.prevLiras)

				watcher.prevRubles = currentRubles
				watcher.prevDollars = currentDollars
				watcher.prevEuros = currentEuros
				watcher.prevLiras = currentLiras
			}
		}
	}()
}

func (watcher *CurrencyRepositoryWatcher) StopWatching() {
	close(watcher.stopChannel)
	watcher.wg.Wait()
}

func logAddedCurrencies[T currency.CurrencyInterface](current, previous []T) {
	if len(current) > len(previous) {
		added := current[len(previous):]
		for _, item := range added {
			fmt.Printf("Added new %s (code=%s, value=%.2f)\n", item.GetName(), item.GetCode(), item.GetValue())
		}
	}
}
