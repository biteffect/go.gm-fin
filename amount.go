package gmfin

import (
	"database/sql/driver"
	"fmt"
	"strconv"
	"strings"
)

/*
 type designed to store money amounts
 with 4 digits after point as internal value
 and work with 2 digits as external value
 example :
    12.0001 -> json -> 12.00

 amount type parsed from json and
 */

type Amount int64

func AmountFromCents(cents int) Amount {
	return Amount(100 * cents)
}

func (m Amount) InCents() int64 {
	return int64(m.BankRound() / 100)
}

func (m Amount) BankRound() Amount {
	abs := m.abs(int64(m))
	cents := abs / 100
	p1 := abs % 100 / 10 // .2346 -> --4-
	p2 := abs % 10       // .2346 -> ---6
	p1 += m.roundDelta(p2, p1%2 > 0)
	cents += m.roundDelta(p1, cents%2 > 0)
	if m < 0 {
		cents = -cents
	}
	return Amount(100 * cents)
}

func (m Amount) Percent(p Amount) Amount {
	return Amount((int64(p) * int64(m)) / 1000000)
}

/*
func (m Amount) PercentOf(p Amount) Amount {
	return Amount((int64(p) * int64(m)) / 1000000)
}
*/
func (m Amount) Add(v Amount) Amount {
	return Amount(int64(m) + int64(v))
}

func (m Amount) Sub(v Amount) Amount {
	return Amount(int64(m) - int64(v))
}

func (m Amount) Negative() Amount {
	return Amount(0 - int64(m))
}

func (m Amount) IntegerPartAsInt() int {
	return int(m / 10000)
}

// Scan implements the sql.Scanner interface for database deserialization.
func (m *Amount) Scan(value interface{}) error {
	if value == nil {
		return m.normalizedUnmarshal("0")
	}
	return m.normalizedUnmarshal(strings.Replace(string(value.([]byte)), ".", "", 1))
}

// Value implements the driver.Valuer interface for database serialization.
func (m Amount) Value() (driver.Value, error) {
	v, err := m.MarshalJSON()
	return string(v), err
}

func (m Amount) String() string {
	if m == 0 {
		return "0"
	}
	s := ""
	if m < 0 {
		s += "-"
	}
	v := m.abs(int64(m))
	s += strconv.FormatInt(v/10000, 10)
	d := v % 10000
	if v%100 > 0 {
		s += "." + fmt.Sprintf("%04d", d)
	} else {
		s += "." + fmt.Sprintf("%02d", d/100)
	}
	return s
}

func (m *Amount) roundDelta(v int64, up bool) int64 {
	if v > 5 || (v == 5 && up) {
		return 1
	}
	return 0
}

func (m *Amount) abs(v int64) int64 {
	if v < 0 {
		return -v
	}
	return v
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (m *Amount) UnmarshalJSON(bytes []byte) error {
	if string(bytes) == "null" {
		return nil
	}
	v := strings.Replace(string(bytes), ",", ".", 1)
	idx := strings.Index(v, ".")
	if idx < 0 {
		v += "0000"
	} else {
		tail := v[idx+1:]
		if len(tail) < 4 {
			tail += strings.Repeat("0", 4-len(tail))
		} else if len(tail) > 4 {
			tail = tail[:4]
		}
		v = v[:idx] + tail
	}

	return m.normalizedUnmarshal(v)
}

// MarshalJSON implements the json.Marshaler interface.
func (m Amount) MarshalJSON() ([]byte, error) {
	return []byte(m.String()), nil
}

func (m *Amount) normalizedUnmarshal(str string) error {
	i, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return err
	}
	*m = Amount(i)
	return nil
}
