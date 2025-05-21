package currency

import (
	"sync"
	"sync/atomic"
)

type IDManager struct {
	counters    map[string]*int32
	mutex       sync.RWMutex
	initialized bool
}

var idManager = &IDManager{
	counters: make(map[string]*int32),
}

func GetNextID(currencyType string) int {
	idManager.mutex.RLock()
	counter, exists := idManager.counters[currencyType]
	idManager.mutex.RUnlock()

	if !exists {
		idManager.mutex.Lock()
		counter, exists = idManager.counters[currencyType]
		if !exists {
			var newCounter int32 = 0
			idManager.counters[currencyType] = &newCounter
			counter = &newCounter
		}
		idManager.mutex.Unlock()
	}

	return int(atomic.AddInt32(counter, 1))
}

func InitializeCounters(maxIDs map[string]int) {
	idManager.mutex.Lock()
	defer idManager.mutex.Unlock()

	if idManager.initialized {
		return
	}

	for currencyType, maxID := range maxIDs {
		var counter int32 = int32(maxID)
		idManager.counters[currencyType] = &counter
	}

	idManager.initialized = true
}

func RegisterCurrencyType(currencyType string, initialID int) {
	idManager.mutex.Lock()
	defer idManager.mutex.Unlock()

	_, exists := idManager.counters[currencyType]
	if !exists {
		var counter int32 = int32(initialID)
		idManager.counters[currencyType] = &counter
	}
}
