package gmfin

var (
	currencies = map[CurrencyCode]*Currency{
		CurrencyUAH: &Currency{CurrencyUAH, 980, "", 2, "\u20b4", "1 $"},
		CurrencyEUR: &Currency{CurrencyEUR, 978, "", 2, "\u20ac", "$1"},
		CurrencyUSD: &Currency{CurrencyUSD, 840, "", 2, "$", "$1"},
		CurrencyRUB: &Currency{CurrencyRUB, 643, "", 2, "\u20bd", "1 $"},
		CurrencyAny: &Currency{CurrencyAny, 0, "", 2, "", "1 $"},
	}
)

func AllCurrencies() []*Currency {
	out := make([]*Currency, 0)
	for _, c := range currencies {
		out = append(out, c)
	}
	return out
}

func SetCurrencies(list []*Currency) {
	nMap := make(map[CurrencyCode]*Currency)
	for _, c := range list {
		nMap[c.Code] = c
	}
	currencies = nMap
}
