package apis

import (
	"encoding/json"
	"fmt"
	"net/url"

	"ucoi4tzlito52agmquc6oopn3zpmr6djz5vvfabtgrhyc6hufpzjtnad.onion/Tochka/tochka-free-market/modules/util"
)

type DOTPayment struct {
	Address string  `json:"address"`
	Percent float64 `json:"percent,omitempty"`
	Amount  float64 `json:"amount,omitempty"`
}

type DOTWalletBalance struct {
	FreeBalance     float64 `json:"free_balance"`
	ReservedBalance float64 `json:"reserved_balance"`
	FrozenBalance   float64 `json:"frozen_balance"`
}

type DOTTransactionResult struct {
	Address string `json:"public_key"`
	Hash    string `json:"hash"`
}

type DOTEscrowMintResult struct {
	PublicKey string `json:"public_key"`
	Hash      string `json:"hash"`
	Id        uint64 `json:"escrow_id"`
}

func GenerateDOTAddress(walletType string) (string, string, error) {
	apiEndpoint := fmt.Sprintf("%s/polkadot/wallets/new", APPLICATION_SETTINGS.PaymentGate)

	response, err := util.DirectPOST(apiEndpoint, url.Values{})
	if err != nil {
		return "", "", err
	}

	var dat map[string]interface{}
	err = json.Unmarshal([]byte(response), &dat)
	if err != nil {
		return "", "", err
	}

	address := dat["public_key"].(string)
	mnemonic := dat["mnemonic"].(string)
	return address, mnemonic, nil
}

func GetAmountOnDOTAddress(address string) (DOTWalletBalance, error) {
	apiEndpoint := fmt.Sprintf("%s/polkadot/wallets/%s", APPLICATION_SETTINGS.PaymentGate, address)

	walletBalance := DOTWalletBalance{}

	body, err := util.DirectGET(apiEndpoint)
	if err != nil {
		return walletBalance, err
	}

	var dat map[string]interface{}
	err = json.Unmarshal([]byte(body), &dat)
	if err != nil {
		return walletBalance, err
	}

	var freeBalance, reservedBalance, frozenBalance float64
	if dat["free_balance"] != nil {
		freeBalance = dat["free_balance"].(float64)
	} else {
		freeBalance = 0.0
	}

	if dat["reserved_balance"] != nil {
		reservedBalance = dat["reserved_balance"].(float64)
	} else {
		reservedBalance = 0.0
	}

	if dat["frozen_balance"] != nil {
		frozenBalance = dat["frozen_balance"].(float64)
	} else {
		frozenBalance = 0.0
	}

	walletBalance.FreeBalance = freeBalance
	walletBalance.ReservedBalance = reservedBalance
	walletBalance.FrozenBalance = frozenBalance

	return walletBalance, nil
}

type DOTPaymentResult struct {
	Hash       string `json:"hash"`
	WalletFrom string `json:"wallet_from"`
	WalletTo   string `json:"wallet_to"`
	Amount     uint64 `json:"amount"`
}

func SendDOT(addressFrom string, dotPayments []DOTPayment) (DOTPaymentResult, error) {
	var (
		paymentsJSON, _ = json.Marshal(dotPayments)
		apiEndpoint     = fmt.Sprintf("%s/polkadot/wallets/%s/send_dot", APPLICATION_SETTINGS.PaymentGate, addressFrom)
		result          = DOTPaymentResult{}
	)

	response, err := util.DirectPOST(apiEndpoint, url.Values{"payments": {string(paymentsJSON)}})
	if err != nil {
		return DOTPaymentResult{}, err
	}

	err = json.Unmarshal([]byte(response), &result)
	return result, err
}

type DOTEscrowMintMetadata struct {
	SellerAddress string  `json:"seller_address"`
	Cid           string  `json:"cid"`
	Value         float64 `json:"value"`
}

func MintEscrow(addressFrom string, dotEscrowMintParams DOTEscrowMintMetadata) (DOTEscrowMintResult, error) {
	var (
		metadataJSON, _ = json.Marshal(dotEscrowMintParams)
		apiEndpoint     = fmt.Sprintf("%s/polkadot/wallets/%s/escrow_mint", APPLICATION_SETTINGS.PaymentGate, addressFrom)
		result          = DOTEscrowMintResult{}
	)

	response, err := util.DirectPOST(apiEndpoint, url.Values{"metadata": {string(metadataJSON)}})
	if err != nil {
		return result, err
	}

	err = json.Unmarshal([]byte(response), &result)
	return result, err
}

func CancelEscrow(escrowId uint64) (DOTTransactionResult, error) {
	var (
		apiEndpoint = fmt.Sprintf("%s/polkadot/escrows/%d/cancel", APPLICATION_SETTINGS.PaymentGate, escrowId)
		result      = DOTTransactionResult{}
	)

	response, err := util.DirectPOST(apiEndpoint, url.Values{})
	if err != nil {
		return result, err
	}

	err = json.Unmarshal([]byte(response), &result)
	return result, err
}

func ReleaseEscrow(escrowId uint64) (DOTTransactionResult, error) {
	var (
		apiEndpoint = fmt.Sprintf("%s/polkadot/escrows/%d/release", APPLICATION_SETTINGS.PaymentGate, escrowId)
		result      = DOTTransactionResult{}
	)

	response, err := util.DirectPOST(apiEndpoint, url.Values{})
	if err != nil {
		return result, err
	}

	err = json.Unmarshal([]byte(response), &result)
	return result, err
}

type PolkadotEscrowCrowdloanStatus struct {
	Status string `json:"status"`
	Time   uint64 `json:"timestamp"`
}

type PolkadotEscrow struct {
	EscrowStatus   PolkadotEscrowCrowdloanStatus `json:"escrow_status"`
	Id             uint64                        `json:"id"`
	MarketplaceId  uint64                        `json:"marketplace_id"`
	Owner          string                        `json:"string"`
	SellerAddress  string                        `json:"seller_address"`
	ShippingStatus PolkadotEscrowCrowdloanStatus `json:"shipping_status"`
	Uri            string                        `json:"uri"`
	Value          uint64                        `json:"value"`
}

func EscrowInfo(escrowId uint64) (PolkadotEscrow, error) {
	var (
		apiEndpoint = fmt.Sprintf("%s/polkadot/escrows/%d", APPLICATION_SETTINGS.PaymentGate, escrowId)
		result      = PolkadotEscrow{}
	)

	response, err := util.DirectGET(apiEndpoint)
	if err != nil {
		return result, err
	}

	err = json.Unmarshal([]byte(response), &result)
	return result, err
}

func EscrowsForBuyer(address string) ([]PolkadotEscrow, error) {
	var (
		apiEndpoint = fmt.Sprintf("%s/polkadot/wallets/%s/escrows", APPLICATION_SETTINGS.PaymentGate, address)
		result      = []PolkadotEscrow{}
	)

	response, err := util.DirectGET(apiEndpoint)
	if err != nil {
		return result, err
	}

	err = json.Unmarshal([]byte(response), &result)
	return result, err
}

func EscrowsForSeller(address string) ([]PolkadotEscrow, error) {
	var (
		apiEndpoint = fmt.Sprintf("%s/polkadot/wallets/%s/deals", APPLICATION_SETTINGS.PaymentGate, address)
		result      = []PolkadotEscrow{}
	)

	response, err := util.DirectGET(apiEndpoint)
	if err != nil {
		return result, err
	}

	err = json.Unmarshal([]byte(response), &result)
	return result, err
}
