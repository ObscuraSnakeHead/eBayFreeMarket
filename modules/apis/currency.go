package apis

import (
	"encoding/json"
	"fmt"
	"time"

	"ucoi4tzlito52agmquc6oopn3zpmr6djz5vvfabtgrhyc6hufpzjtnad.onion/Tochka/tochka-free-market/modules/util"
)

func GetCurrencyRates(baseCurrency string) (map[string]float64, error) {
	rates := map[string]float64{}

	apiEndpoint := fmt.Sprintf("%s/currency/%s", APPLICATION_SETTINGS.PaymentGate, baseCurrency)
	response, err := util.DirectGET(apiEndpoint)
	if err != nil {
		return rates, err
	}

	err = json.Unmarshal([]byte(response), &rates)
	return rates, err
}

func GetHistoricCurrencyRate(baseCurrency string, dt time.Time) float64 {
	// 2019-11-25 price
	if baseCurrency == "ethereum" {
		return 148.5
	}
	if baseCurrency == "bitcoin" {
		return 7220.80
	}

	return 1.0
}
