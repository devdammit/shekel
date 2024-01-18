package currency

import xcurrency "golang.org/x/text/currency"

var (
	THB = Code{xcurrency.THB}
	RUB = Code{xcurrency.RUB}
	USD = Code{xcurrency.USD}
	EUR = Code{xcurrency.EUR}
)
