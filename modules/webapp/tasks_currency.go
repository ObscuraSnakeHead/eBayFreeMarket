package webapp

import (
	"github.com/jasonlvhit/gocron"
)

func StartCurrencyCron() {
	UpdateCurrencyRates()
	gocron.Every(60).Minute().Do(UpdateCurrencyRates)
}
