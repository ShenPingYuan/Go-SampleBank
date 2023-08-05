package util

const (
	CNY = "CNY"
	EUR = "EUR"
	USD = "USD"
)

func IsValidCurrency(currency string) bool {
	switch currency {
	case CNY, EUR, USD:
		return true
	}
	return false
}