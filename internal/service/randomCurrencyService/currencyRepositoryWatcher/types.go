package currencyRepositoryWatcher

import (
	"OtusGo/internal/model/currency"
	"OtusGo/internal/repository/currencyRepository"
	"fmt"
	"time"
)

type CurrencyRepositoryWatcher struct {
	currencyRepository *currencyRepository.CurrencyRepository
	prevRubles         []*currency.Ruble
	prevDollars        []*currency.Dollar
	prevEuros          []*currency.Euro
	prevLiras          []*currency.Lira
}

func NewCurrencyRepositoryWatcher(currencyRepository *currencyRepository.CurrencyRepository) *CurrencyRepositoryWatcher {
	watcher := &CurrencyRepositoryWatcher{
		currencyRepository: currencyRepository,
	}

	go watcher.startWatching()

	return watcher
}

func (watcher *CurrencyRepositoryWatcher) startWatching() {
	ticker := time.NewTicker(200 * time.Millisecond)
	defer ticker.Stop()

	for range ticker.C {
		currentRubles := watcher.currencyRepository.GetRubles()
		currentDollars := watcher.currencyRepository.GetDollars()
		currentEuros := watcher.currencyRepository.GetEuros()
		currentLiras := watcher.currencyRepository.GetLiras()

		if len(currentRubles) > len(watcher.prevRubles) {
			added := currentRubles[len(watcher.prevRubles):]
			for _, ruble := range added {
				fmt.Printf("Added Ruble: Name=%s, Code=%s, Value=%.2f\n", ruble.GetName(), ruble.GetCode(), ruble.GetValue())
			}
		}

		if len(currentDollars) > len(watcher.prevDollars) {
			added := currentDollars[len(watcher.prevDollars):]
			for _, dollar := range added {
				fmt.Printf("Added Dollar: Name=%s, Code=%s, Value=%.2f\n", dollar.GetName(), dollar.GetCode(), dollar.GetValue())
			}
		}

		if len(currentEuros) > len(watcher.prevEuros) {
			added := currentEuros[len(watcher.prevEuros):]
			for _, euro := range added {
				fmt.Printf("Added Euro: Name=%s, Code=%s, Value=%.2f\n", euro.GetName(), euro.GetCode(), euro.GetValue())
			}
		}

		if len(currentLiras) > len(watcher.prevLiras) {
			added := currentLiras[len(watcher.prevLiras):]
			for _, lira := range added {
				fmt.Printf("Added Lira: Name=%s, Code=%s, Value=%.2f\n", lira.GetName(), lira.GetCode(), lira.GetValue())
			}
		}

		watcher.prevRubles = currentRubles
		watcher.prevDollars = currentDollars
		watcher.prevEuros = currentEuros
		watcher.prevLiras = currentLiras
	}
}
