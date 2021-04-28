package gmfin

var currencyStore ICurrencyStore

func SetCurrencyStore(store ICurrencyStore) {
	currencyStore = store
}

func setDummyCurrencyStore() {
	// currencyStore = store
}

type uaCurrencyStore struct {
}
