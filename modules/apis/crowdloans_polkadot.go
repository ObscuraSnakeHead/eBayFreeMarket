package apis

import (
	"encoding/json"
	"fmt"
	"net/url"

	"ucoi4tzlito52agmquc6oopn3zpmr6djz5vvfabtgrhyc6hufpzjtnad.onion/Tochka/tochka-free-market/modules/util"
)

type DOTCrowdloanMintMetadata struct {
	GoalValue      float64 `json:"goal_value"`
	GoalDate       uint64  `json:"goal_date"`
	WeeklyInterest uint16  `json:"weekly_interest"`
}

type DOTCrowdloanMintResult struct {
	PublicKey string `json:"public_key"`
	Hash      string `json:"hash"`
	Id        uint64 `json:"crowdloan_id"`
}

type DOTFundPaybackCrowdloanMetadata struct {
	CrowdloanId uint64  `json:"crowdloan_id"`
	Value       float64 `json:"value"`
}

func MintCrowdloan(
	addressFrom string,
	dotCrowdloanMintParams DOTCrowdloanMintMetadata,
) (DOTCrowdloanMintResult, error) {
	var (
		metadataJSON, _ = json.Marshal(dotCrowdloanMintParams)
		apiEndpoint     = fmt.Sprintf("%s/polkadot/wallets/%s/crowdloan_mint", APPLICATION_SETTINGS.PaymentGate, addressFrom)
		result          = DOTCrowdloanMintResult{}
	)

	response, err := util.DirectPOST(apiEndpoint, url.Values{"metadata": {string(metadataJSON)}})
	if err != nil {
		return result, err
	}

	err = json.Unmarshal([]byte(response), &result)
	return result, err
}

type PolkdadotCrowdloan struct {
	CloseDate       *uint                         `json:"close_date"`
	ClosingAmount   float64                       `json:"closing_amount"`
	Value           float64                       `json:"value"`
	EscrowStatus    PolkadotEscrowCrowdloanStatus `json:"status"`
	GoalValue       float64                       `json:"goal_value"`
	GoalDate        uint64                        `json:"goal_date"`
	Id              uint64                        `json:"id"`
	NumberOfLenders uint                          `json:"n_lenders"`
	Owner           string                        `json:"owner"`
	WeeklyIterest   uint64                        `json:"weekly_interest"`
}

func CrowdloanInfo(loanId uint64) (PolkdadotCrowdloan, error) {
	var (
		apiEndpoint = fmt.Sprintf("%s/polkadot/loans/%d", APPLICATION_SETTINGS.PaymentGate, loanId)
		result      = PolkdadotCrowdloan{}
	)

	response, err := util.DirectGET(apiEndpoint)
	if err != nil {
		return result, err
	}

	err = json.Unmarshal([]byte(response), &result)
	return result, err
}

func FundCrowdloan(addressFrom string, dotCrowdloanMintParams DOTFundPaybackCrowdloanMetadata) (DOTTransactionResult, error) {
	var (
		metadataJSON, _ = json.Marshal(dotCrowdloanMintParams)
		apiEndpoint     = fmt.Sprintf("%s/polkadot/wallets/%s/crowdloan_fund", APPLICATION_SETTINGS.PaymentGate, addressFrom)
		result          = DOTTransactionResult{}
	)

	response, err := util.DirectPOST(apiEndpoint, url.Values{"metadata": {string(metadataJSON)}})
	if err != nil {
		return result, err
	}

	err = json.Unmarshal([]byte(response), &result)
	return result, err
}

func CrowdloandMintedByAddress(address string) ([]PolkdadotCrowdloan, error) {
	var (
		apiEndpoint = fmt.Sprintf("%s/polkadot/wallets/%s/loans", APPLICATION_SETTINGS.PaymentGate, address)
		result      = []PolkdadotCrowdloan{}
	)

	response, err := util.DirectGET(apiEndpoint)
	if err != nil {
		return result, err
	}

	err = json.Unmarshal([]byte(response), &result)
	return result, err
}

type PolkdadotCrowdloanLend struct {
	Id          uint64             `json:"id"`
	CrowdloanId uint64             `json:"crowdloan_id"`
	LendValue   uint64             `json:"lend_value"`
	Index       uint64             `json:"index"`
	Loan        PolkdadotCrowdloan `json:"loan"`
}

func CrowdloandFundedByAddress(address string) ([]PolkdadotCrowdloanLend, error) {
	var (
		apiEndpoint = fmt.Sprintf("%s/polkadot/wallets/%s/lends", APPLICATION_SETTINGS.PaymentGate, address)
		result      = []PolkdadotCrowdloanLend{}
	)

	response, err := util.DirectGET(apiEndpoint)
	if err != nil {
		return result, err
	}

	err = json.Unmarshal([]byte(response), &result)
	return result, err
}

type DOTWithdawCrowdloanMetadata struct {
	CrowdloanId uint64 `json:"crowdloan_id"`
}

func WithdrawCrowdloan(addressFrom string, dotCrowdloanMintParams DOTWithdawCrowdloanMetadata) (DOTTransactionResult, error) {
	var (
		metadataJSON, _ = json.Marshal(dotCrowdloanMintParams)
		apiEndpoint     = fmt.Sprintf("%s/polkadot/wallets/%s/crowdloan_payback", APPLICATION_SETTINGS.PaymentGate, addressFrom)
		result          = DOTTransactionResult{}
	)

	response, err := util.DirectPOST(apiEndpoint, url.Values{"metadata": {string(metadataJSON)}})
	if err != nil {
		return result, err
	}

	err = json.Unmarshal([]byte(response), &result)
	return result, err
}

type DOTPaybackCrowdloanMetadata struct {
	CrowdloanId uint64 `json:"crowdloan_id"`
}

func PaybackCrowdloan(addressFrom string, dotCrowdloanMintParams DOTPaybackCrowdloanMetadata) (DOTTransactionResult, error) {
	var (
		metadataJSON, _ = json.Marshal(dotCrowdloanMintParams)
		apiEndpoint     = fmt.Sprintf("%s/polkadot/wallets/%s/crowdloan_payback", APPLICATION_SETTINGS.PaymentGate, addressFrom)
		result          = DOTTransactionResult{}
	)

	response, err := util.DirectPOST(apiEndpoint, url.Values{"metadata": {string(metadataJSON)}})
	if err != nil {
		return result, err
	}

	err = json.Unmarshal([]byte(response), &result)
	return result, err
}

type DOTPayoutCrowdloanMetadata struct {
	CrowdloanId uint64 `json:"crowdloan_id"`
}

func CollectCrowdloanPayout(addressFrom string, dotCrowdloanMintParams DOTPayoutCrowdloanMetadata) (DOTTransactionResult, error) {
	var (
		metadataJSON, _ = json.Marshal(dotCrowdloanMintParams)
		apiEndpoint     = fmt.Sprintf("%s/polkadot/wallets/%s/crowdloan_payout", APPLICATION_SETTINGS.PaymentGate, addressFrom)
		result          = DOTTransactionResult{}
	)

	response, err := util.DirectPOST(apiEndpoint, url.Values{"metadata": {string(metadataJSON)}})
	if err != nil {
		return result, err
	}

	err = json.Unmarshal([]byte(response), &result)
	return result, err
}
