package gmfin

import (
	"strings"
)

type CurrencyAmount struct {
	Amount   Amount    `pg:",type:numeric(20,4),notnull,default:0"`
	Currency *Currency `json:",omitempty" pg:",type:varchar(3),notnull,default:'UAH'"`
}

func (a CurrencyAmount) InCents() int64 {
	return a.Amount.InCents()
}

func (a CurrencyAmount) String() string {
	return strings.Replace(
		strings.Replace(a.Currency.Template, "1", a.Amount.String(), 1),
		"$", a.Currency.Grapheme, 1)
}
