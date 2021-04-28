package gmfin

import (
	"database/sql/driver"
	"fmt"
	"strings"
)

type CurrencyCode string

type ICurrencyStore interface {
	ByCode(CurrencyCode) (*Currency, error)
	ByNumericCode(int) (*Currency, error)
}

type currencyBase struct {
	Code        CurrencyCode `pg:",notnull,pk,type:varchar(3)"`
	NumericCode int          `json:"-" pg:"numeric,notnull"`
	Name        string       `pg:",notnull,type:varchar(24)"`
	Fraction    int          `pg:",notnull"`
	Grapheme    string       `pg:",notnull"`
	Template    string       `pg:",notnull"`
}

type Currency struct {
	currencyBase
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
	if currencyStore == nil {
		setDummyCurrencyStore()
	}
	v, err := currencyStore.ByCode(CurrencyCode(str))
	if err != nil {
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
		if v, err := currencyStore.ByCode(CurrencyCode(v.(string))); err == nil {
			return v, nil
		}
	case int:
		if v, err := currencyStore.ByNumericCode(v.(int)); err == nil {
			return v, nil
		}
	default:
	}
	return nil, fmt.Errorf("unsupported currency: %v", v)
}
