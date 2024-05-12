/*
* Created on 09 May 2024
* @author Sai Sumanth
 */

package utils

// constants for all supported currencies
const (
	USD = "USD"
	INR = "INR"
	EUR = "EUR"
)

// IsSupportedCurrency will return True if the currency is supported
func IsSupportedCurrency(currency string) bool {
	switch currency {
	case USD, INR, EUR:
		return true
	}
	return false
}