package gmfin

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type CardNumberType uint64
type CardIssuerType string

const (
	Visa       CardIssuerType = "visa"
	MasterCard CardIssuerType = "mastercard"
	Amex       CardIssuerType = "american express"
	Diners     CardIssuerType = "diners"
	Discover   CardIssuerType = "discover"
	JCB        CardIssuerType = "jcb"
	Prostir    CardIssuerType = "prostir"
	Other      CardIssuerType = "other"
)

type CreditCard struct {
	CardNumber       uint64     `json:"cardNumber" pg:"cardNumber,notnull"`
	CardSecurityCode uint       `json:"cardSecurityCode"`
	CardholderName   string     `json:"cardholderName"`
	ExpiryMonth      time.Month `json:"expiryMonth"`
	ExpiryYear       uint       `json:"expiryYear"`
}

func (c *CreditCard) NumberString() string {
	return strconv.FormatUint(c.CardNumber, 10)
}

func (c *CreditCard) Mask() string {
	str := c.NumberString()
	return str[:4] + strings.Repeat("*", len(str)-8) + str[len(str)-4:]
}

func (c *CreditCard) MaskFull() string {
	str := c.NumberString()
	if len(str) < 11 {
		return str
	}
	return str[:6] + strings.Repeat("*", len(str)-10) + str[len(str)-4:]
}

func (c *CreditCard) SecurityCodeString() string {
	return strconv.Itoa(int(c.CardSecurityCode))
}

func (c *CreditCard) Issuer() CardIssuerType {
	regVisa, _ := regexp.Compile(`^4[0-9]{12}(?:[0-9]{3})?$`)
	regMaster, _ := regexp.Compile(`^5[1-5][0-9]{14}$`)
	regAmex, _ := regexp.Compile(`^3[47][0-9]{13}$`)
	regDiners, _ := regexp.Compile(`^3(?:0[0-5]|[68][0-9])[0-9]{11}$`)
	regDiscover, _ := regexp.Compile(`^6(?:011|5[0-9]{2})[0-9]{12}$`)
	regJCB, _ := regexp.Compile(`^(?:2131|1800|35\d{3})\d{11}$`)
	regProstir, _ := regexp.Compile(`^9[0-9]{12}(?:[0-9]{3})?$`)
	reg := map[CardIssuerType]interface{}{
		Visa:       regVisa,
		MasterCard: regMaster,
		Amex:       regAmex,
		Diners:     regDiners,
		Discover:   regDiscover,
		JCB:        regJCB,
		Prostir:    regProstir,
	}
	for t, r := range reg {
		if r.(*regexp.Regexp).MatchString(c.NumberString()) {
			return t
		}
	}
	return Other
}

func (c *CreditCard) Validate() error {
	if c.ExpiryMonth < 1 || c.ExpiryMonth > 12 {
		return fmt.Errorf("invalid expiry month: %v", c.CardSecurityCode)
	}
	if c.ExpiryYear < 1 || c.ExpiryMonth > 12 {
		return fmt.Errorf("invalid expiry month: %v", c.CardSecurityCode)
	}
	now := time.Now()
	exp := time.Date(int(c.ExpiryYear), c.ExpiryMonth,
		1, 0, 0, 0, 0, now.Location())
	// Validate: card expired
	if exp.Before(now) {
		return fmt.Errorf("card expired: %v", exp.Format("01/2006"))
	}
	// Validate: cvv2
	if c.CardSecurityCode < 1 || c.CardSecurityCode > 999 {
		return fmt.Errorf("card CVV2 length invalid: %v", c.CardSecurityCode)
	}
	return c.ValidateNumber()
}

func (c *CreditCard) ValidateNumber() error {

	// Validate: card length ( 13 ... 19 digits )
	if c.CardNumber < 1000000000000 || c.CardNumber >= 10000000000000000000 {
		return fmt.Errorf("card number length invalid: %v", c.CardNumber)
	}
	// Validate: Luhn algorithm
	if !Luhn(c.NumberString()) {
		return fmt.Errorf("card number failed the luhn algorithm check: %v", c.CardNumber)
	}
	return nil
}

func (c *CreditCard) IsTestCard() bool {
	switch c.CardNumber {
	case 4242424242424242,
		4012888888881881,
		4000056655665556,
		5555555555554444,
		5200828282828210,
		5105105105105100,
		378282246310005,
		371449635398431,
		6011111111111117,
		6011000990139424,
		30569309025904,
		38520000023237,
		3530111333300000,
		3566002020360505:
		return true
	}
	return false
}

func Luhn(number string) bool {
	l := len(number)
	if l < 13 || l > 16 {
		return false
	}
	parity := l % 2
	sum, err := strconv.Atoi(number[l-1:])
	if err != nil {
		return false
	}
	for i := 0; i < l-1; i++ {
		num, err := strconv.Atoi(string(number[i]))
		if err != nil {
			return false
		}
		if i%2 == parity {
			num = num * 2
			if num > 9 {
				num = num - 9
			}
		}
		sum += num
	}
	return sum%10 == 0
}
