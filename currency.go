package gmfin

import (
	"database/sql/driver"
	"fmt"
	"strings"
)

type CurrencyCode string

const (
	CurrencyAny CurrencyCode = "ANY"
	CurrencyUAH CurrencyCode = "UAH"
	CurrencyRUB CurrencyCode = "RUB"
	CurrencyUSD CurrencyCode = "USD"
	CurrencyEUR CurrencyCode = "EUR"
)

type Currency struct {
	Code        CurrencyCode `pg:",notnull,pk,type:varchar(3)"`
	NumericCode int          `json:"-" pg:"numeric,notnull"`
	Name        string       `pg:",notnull,type:varchar(24)"`
	Fraction    int          `pg:",notnull"`
	Grapheme    string       `pg:",notnull"`
	Template    string       `pg:",notnull"`
}

func (c *Currency) String() string {
	return string(c.Code)
}

func (c *Currency) Equals(cc *Currency) bool {
	return c.NumericCode == cc.NumericCode
}

// Scan implements the sql.Scanner interface for database deserialization.
func (c *Currency) Scan(value interface{}) error {
	return c.UnmarshalJSON(value.([]byte))
}

// Value implements the driver.Valuer interface for database serialization.
func (c Currency) Value() (driver.Value, error) {
	// TODO check value
	return c.Code, nil
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (c *Currency) UnmarshalJSON(bytes []byte) error {
	str := strings.ToUpper(strings.Trim(string(bytes), "\" "))
	if len(str) != 3 {
		return fmt.Errorf("unsupported currency: %v", str)
	}
	v, ok := currencies[CurrencyCode(str)]
	if !ok {
		return fmt.Errorf("unsupported currency: %v", str)
	}
	*c = *v
	return nil
}

// MarshalJSON implements the json.Marshaler interface.
func (c Currency) MarshalJSON() ([]byte, error) {
	return []byte("\"" + c.Code + "\""), nil
}

func NewCurrency(v interface{}) (*Currency, error) {
	switch v.(type) {
	case string:
		if v, ok := currencies[CurrencyCode(v.(string))]; ok {
			return v, nil
		}
	case int:
		for _, c := range currencies {
			if c.NumericCode == v.(int) {
				return c, nil
			}
		}
	default:
	}
	return nil, fmt.Errorf("unsupported currency: %v", v)
}
