package apis

import (
	"encoding/json"
	"fmt"
	"net/url"

	"ucoi4tzlito52agmquc6oopn3zpmr6djz5vvfabtgrhyc6hufpzjtnad.onion/Tochka/tochka-free-market/modules/util"
)

type XMRPayment struct {
	Address string  `json:"address"`
	Percent float64 `json:"percent,omitempty"`
	Amount  float64 `json:"amount,omitempty"`
}

type XMRWalletBalance struct {
	Balance         float64 `json:"balance"`
	UnlockedBalance float64 `json:"unlocked_balance"`
}

type XMRPaymentResult struct {
	Amount        int64  `json:"amount"`
	Fee           int64  `json:"fee"`
	MultisigTxset string `json:"multisig_txset"`
	TxBlob        string `json:"tx_blob"`
	TxHash        string `json:"tx_hash"`
	TxMetadata    string `json:"tx_metadata"`
	TxKey         string `json:"tx_key"`
	UnsignedTxset string `json:"unsigned_txset"`
}

func GenerateXMRAddress(walletType string) (string, error) {
	apiEndpoint := fmt.Sprintf("%s/monero/wallets/new", APPLICATION_SETTINGS.PaymentGate)

	response, err := util.DirectPOST(apiEndpoint, url.Values{"type": {walletType}})
	if err != nil {
		return "", err
	}

	var dat map[string]interface{}
	err = json.Unmarshal([]byte(response), &dat)
	if err != nil {
		return "", err
	}

	address := dat["address"].(string)
	return address, nil
}

func GetAmountOnXMRAddress(address string) (XMRWalletBalance, error) {
	apiEndpoint := fmt.Sprintf("%s/monero/wallets/%s", APPLICATION_SETTINGS.PaymentGate, address)

	walletBalance := XMRWalletBalance{}

	body, err := util.DirectGET(apiEndpoint)
	if err != nil {
		return walletBalance, err
	}

	var dat map[string]interface{}
	err = json.Unmarshal([]byte(body), &dat)
	if err != nil {
		return walletBalance, err
	}

	var balance, unlockedBalance float64
	if dat["balance"] != nil {
		balance = dat["balance"].(float64)
	} else {
		balance = 0.0
	}

	if dat["unlocked_balance"] != nil {
		unlockedBalance = dat["unlocked_balance"].(float64)
	} else {
		unlockedBalance = 0.0
	}

	walletBalance.Balance = balance
	walletBalance.UnlockedBalance = unlockedBalance

	return walletBalance, nil
}

func SendXMR(addressFrom string, xmrPayements []XMRPayment) ([]XMRPaymentResult, error) {
	var (
		paymentsJSON, _ = json.Marshal(xmrPayements)
		apiEndpoint     = fmt.Sprintf("%s/monero/wallets/%s/send", APPLICATION_SETTINGS.PaymentGate, addressFrom)
		result          = []XMRPaymentResult{}
	)

	response, err := util.DirectPOST(apiEndpoint, url.Values{"payments": {string(paymentsJSON)}})
	if err != nil {
		return result, err
	}

	println(
		fmt.Sprintf(
			"Sending request to payaka:\ncurl %s -d `payments=%s`",
			apiEndpoint,
			string(paymentsJSON),
		))

	err = json.Unmarshal([]byte(response), &result)
	return result, err
}
