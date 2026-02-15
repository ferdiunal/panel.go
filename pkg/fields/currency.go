package fields

// Currency, money field'larında kullanılacak ISO 4217 currency kodunu temsil eder.
type Currency string

const (
	CurrencyUSD Currency = "USD"
	CurrencyEUR Currency = "EUR"
	CurrencyTRY Currency = "TRY"
	CurrencyGBP Currency = "GBP"
	CurrencyJPY Currency = "JPY"
	CurrencyCHF Currency = "CHF"
	CurrencyCAD Currency = "CAD"
	CurrencyAUD Currency = "AUD"
	CurrencyCNY Currency = "CNY"
)

// DefaultCurrencies, money field için önerilen başlangıç para birimi listesidir.
var DefaultCurrencies = []Currency{
	CurrencyUSD,
	CurrencyEUR,
	CurrencyTRY,
	CurrencyGBP,
}

func currencyCodes(currencies []Currency) []string {
	codes := make([]string, 0, len(currencies))
	for _, currency := range currencies {
		codes = append(codes, string(currency))
	}
	return codes
}
