package gopartpicker

import (
	"strconv"
	"strings"
	"unicode"
)

// Represents the total price of an item as well as additional fees.
type Price struct {
	// The base price of the item, without shipping, discounts or tax.
	Base float64
	// The price of shipping for the item.
	Shipping float64
	// The price of tax for the item.
	Tax float64
	// The price of discounts for the item.
	Discounts float64
	// The total price of the item.
	Total float64
	// The currency of the price of the item, e.g. £, $.
	Currency string
	// A string representing the item's total price, e.g. $1000 or 5500 RON.
	TotalString string
}

// Converts a string representation of a into a float and a string representing the currency.
func StringPriceToFloat(price string) (float64, string, error) {
	price = strings.TrimSpace(price)

	if price == "" {
		return 0, "", nil
	}

	var currency string
	var number string

	for _, char := range price {
		if char == ' ' || char == '+' {
			continue
		} else if char == '.' || char == ',' {
			number += "."
		} else if unicode.IsDigit(char) {
			number += string(char)
		} else {
			currency += string(char)
		}
	}

	float, err := strconv.ParseFloat(number, 64)

	if err != nil {
		return 0, "", err
	}

	return float, currency, nil
}
