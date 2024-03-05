package util

// Constants for supported currencies
const (
	USD = "USD"
	EUR = "EUR"
	CAD = "CAD"
	VND = "VND"
)

// IsSupportedCurrency checks if the currency is supported
func IsSupportedCurrency(currency string) bool {
	switch currency {
	case USD, EUR, CAD, VND:
		return true
	}
	return false
}
