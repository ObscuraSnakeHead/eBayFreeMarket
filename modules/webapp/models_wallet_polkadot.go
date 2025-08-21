package webapp

import (
	"errors"
	"fmt"
	"time"

	"ucoi4tzlito52agmquc6oopn3zpmr6djz5vvfabtgrhyc6hufpzjtnad.onion/Tochka/tochka-free-market/modules/apis"
)

/*
	Models
*/

type UserPolkadotWallet struct {
	PublicKey string `json:"public_key" gorm:"primary_key"`
	Mnemonic  string `json:"mnemonic"`
	UserUuid  string `json:"-" gorm:"index"`
	IsLocked  bool   `json:"-"`

	User User `json:"-"`

	CreatedAt *time.Time `json:"created_at,omitempty"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}

type UserPolkadotWalletBalance struct {
	ID              int     `json:"id" gorm:"primary_key"`
	PublicKey       string  `json:"public_key" gorm:"index"`
	FreeBalance     float64 `json:"free_balance"`
	ReservedBalance float64 `json:"reserved_balance"`
	FrozenBalance   float64 `json:"frozen_balance"`
	Type            string  `json:"type"`

	CreatedAt          time.Time          `json:"created_at"`
	UserPolkadotWallet UserPolkadotWallet `gorm:"AssociationForeignKey:PublicKey" json:"-"`
}

type UserPolkadotWallets []UserPolkadotWallet

type UserPolkadotWalletAction struct {
	ID                 int    `json:"id" gorm:"primary_key"`
	UserUuid           string `json:"-" gorm:"user_uuid"`
	PaymentReceiptUuid string `json:"-" gorm:"payment_receipt_uuid"`
	PublicKey          string `json:"-" gorm:"index"`

	Action string  `json:"action"`
	Amount float64 `json:"amount"`

	CreatedAt time.Time `json:"created_at"`

	User               User               `json:"-"`
	UserPolkadotWallet UserPolkadotWallet `gorm:"AssociationForeignKey:PublicKey" json:"-"`
	PaymentReceipt     PaymentReceipt     `json:"payment_receipt"`
}

/*
	Model Methods
*/

func (w UserPolkadotWallet) UpdateBalance(force bool) (UserPolkadotWalletBalance, error) {
	d, _ := time.ParseDuration("1m")
	now := time.Now()

	if w.UpdatedAt != nil && !w.UpdatedAt.Add(d).Before(now) && !force {
		return UserPolkadotWalletBalance{}, errors.New("polkadot wallet was updated recently. please wait")
	}

	w.UpdatedAt = &now

	err := w.Save()
	if err != nil {
		return UserPolkadotWalletBalance{}, err
	}

	balance, err := apis.GetAmountOnDOTAddress(w.PublicKey)
	if err != nil {
		return UserPolkadotWalletBalance{}, err
	}

	currentBalance := w.CurrentBalance()
	if currentBalance.FreeBalance != balance.FreeBalance || currentBalance.FrozenBalance != balance.FrozenBalance || currentBalance.ReservedBalance != balance.ReservedBalance {

		uwb := UserPolkadotWalletBalance{
			PublicKey:       w.PublicKey,
			FreeBalance:     balance.FreeBalance,
			FrozenBalance:   balance.FrozenBalance,
			ReservedBalance: balance.ReservedBalance,
			CreatedAt:       now,
			Type:            "blockchain",
		}
		uwb.Save()

		var (
			action string
			amount float64
		)

		if currentBalance.FreeBalance != balance.FreeBalance {
			action = "free balance updated"
			amount = balance.FreeBalance
		} else if currentBalance.FrozenBalance != balance.FrozenBalance {
			action = "frozen balance updated"
			amount = balance.FrozenBalance
		} else if currentBalance.ReservedBalance != balance.ReservedBalance {
			action = "reserved balance updated"
			amount = balance.FrozenBalance
		}

		uwa := UserPolkadotWalletAction{
			UserUuid:  w.UserUuid,
			PublicKey: w.PublicKey,
			Action:    action,
			Amount:    amount,
			CreatedAt: time.Now(),
		}
		uwa.Save()

		return uwb, nil
	}

	return currentBalance, nil
}

func (w UserPolkadotWallet) CurrentBalance() UserPolkadotWalletBalance {
	var uwb UserPolkadotWalletBalance

	database.
		Where(&UserPolkadotWalletBalance{PublicKey: w.PublicKey}).
		Order("created_at DESC").
		First(&uwb)

	return uwb
}

/*
	Database methods
*/

func (w UserPolkadotWallet) Validate() error {
	if w.PublicKey == "" {
		return errors.New("wrong wallet")
	}
	if w.Mnemonic == "" {
		return errors.New("wrong mnemonic")
	}
	return nil
}

func (w UserPolkadotWallet) Save() error {
	err := w.Validate()
	if err != nil {
		return err
	}
	return w.SaveToDatabase()
}

func (w UserPolkadotWallet) Remove() error {
	return database.Delete(w).Error
}

func (w UserPolkadotWallet) SaveToDatabase() error {
	if existing, _ := FindUserPolkadotWalletByPublicKey(w.PublicKey); existing == nil {
		return database.Create(&w).Error
	}
	return database.Save(&w).Error
}

// UserPolkadotWalletBalance

func (w UserPolkadotWalletBalance) Validate() error {
	if w.PublicKey == "" {
		return errors.New("wrong wallet")
	}
	return nil
}

func (w UserPolkadotWalletBalance) Save() error {
	err := w.Validate()
	if err != nil {
		return err
	}
	return w.SaveToDatabase()
}

func (w UserPolkadotWalletBalance) Remove() error {
	return database.Delete(w).Error
}

func (w UserPolkadotWalletBalance) SaveToDatabase() error {
	if w.ID != 0 {
		return database.Save(&w).Error
	}
	w.ID = GetNextUserPolkadotWalletbalanceID()
	return database.Create(&w).Error
}

func (w UserPolkadotWalletAction) Validate() error {
	if w.UserUuid == "" {
		return errors.New("wrong UserUuid")
	}
	return nil
}

func (w UserPolkadotWalletAction) Save() error {
	err := w.Validate()
	if err != nil {
		return err
	}
	return w.SaveToDatabase()
}

func (w UserPolkadotWalletAction) Remove() error {
	return database.Delete(w).Error
}

func (w UserPolkadotWalletAction) SaveToDatabase() error {
	if w.ID != 0 {
		return database.Save(&w).Error
	}
	w.ID = GetNextUserPolkadotWalletActionID()
	return database.Create(&w).Error
}

/*
	Model Methods
*/

func (umws UserPolkadotWallets) Balance() apis.DOTWalletBalance {
	var balance apis.DOTWalletBalance

	for _, umw := range umws {
		balance.FreeBalance += umw.CurrentBalance().FreeBalance
		balance.FrozenBalance += umw.CurrentBalance().FrozenBalance
		balance.ReservedBalance += umw.CurrentBalance().ReservedBalance
	}

	return balance
}

func (umw UserPolkadotWallet) Balance() apis.DOTWalletBalance {
	var balance apis.DOTWalletBalance

	balance.FreeBalance = umw.CurrentBalance().FreeBalance
	balance.FrozenBalance = umw.CurrentBalance().FrozenBalance
	balance.ReservedBalance = umw.CurrentBalance().ReservedBalance

	return balance
}

func (w UserPolkadotWallet) prepareTransaction(address string, amount float64) (apis.DOTPayment, error) {
	if amount > w.Balance().FreeBalance {
		return apis.DOTPayment{}, errors.New("amount is greater than free balance")
	}

	if !polkadotRegexp.MatchString(address) {
		return apis.DOTPayment{}, errors.New("wrong DOT address")
	}

	if w.IsLocked {
		return apis.DOTPayment{}, errors.New("wallet is currently locked")
	}

	return apis.DOTPayment{Address: address, Amount: amount}, nil
}

func (umw UserPolkadotWallet) Send(address string, amount float64) (PaymentReceipt, error) {
	payment, err := umw.prepareTransaction(address, amount)
	if err != nil {
		return PaymentReceipt{}, err
	}

	transactionResult, err := apis.SendDOT(umw.PublicKey, []apis.DOTPayment{payment})
	if err != nil {
		return PaymentReceipt{}, err
	}

	receipt, err := CreateDOTTransactionReceiptForPayment(transactionResult)
	if err != nil {
		return PaymentReceipt{}, err
	}

	uwa := UserPolkadotWalletAction{
		Action:             fmt.Sprintf("Sent %f DOT to %s", amount, address),
		UserUuid:           umw.UserUuid,
		CreatedAt:          time.Now(),
		PaymentReceiptUuid: receipt.Uuid,
	}

	return receipt, uwa.Save()
}

/*
	Crowdloan Methods
*/

func (umw UserPolkadotWallet) MintCrowdloan(
	goalValue float64,
	goalDate time.Time,
	weeklyInterest uint16,
) (Crowdloan, error) {
	mintResult, err := apis.MintCrowdloan(umw.PublicKey, apis.DOTCrowdloanMintMetadata{
		GoalValue:      goalValue,
		GoalDate:       uint64(goalDate.Unix()),
		WeeklyInterest: weeklyInterest,
	})
	if err != nil {
		return Crowdloan{}, err
	}

	receipt, err := CreateDOTTransactionReceiptForCrowdloanMint(mintResult)
	if err != nil {
		return Crowdloan{}, err
	}

	uwa := UserPolkadotWalletAction{
		Action:             fmt.Sprintf("Minted Crowdloan (%f DOT)", goalValue),
		UserUuid:           umw.UserUuid,
		CreatedAt:          time.Now(),
		PaymentReceiptUuid: receipt.Uuid,
	}

	err = uwa.Save()
	if err != nil {
		return Crowdloan{}, err
	}

	return CreateCrowdloan(
		fmt.Sprintf("polkadot-%d", mintResult.Id),
		"polkadot",
		goalValue,
		goalDate,
		uint16(weeklyInterest),
		umw.UserUuid,
		mintResult.Id,
	)
}

func (umw UserPolkadotWallet) MintedCrowdloans() ([]apis.PolkdadotCrowdloan, error) {
	// TODO: sync with database
	return apis.CrowdloandMintedByAddress(umw.PublicKey)
}

func (umw UserPolkadotWallet) FundedCrowdloans() ([]apis.PolkdadotCrowdloanLend, error) {
	// TODO: sync with database
	return apis.CrowdloandFundedByAddress(umw.PublicKey)
}

func (umw UserPolkadotWallet) FundCrowdloan(crowdloanId uint64, value float64) (PaymentReceipt, error) {
	fundResult, err := apis.FundCrowdloan(umw.PublicKey, apis.DOTFundPaybackCrowdloanMetadata{
		CrowdloanId: crowdloanId,
		Value:       value,
	})
	if err != nil {
		return PaymentReceipt{}, err
	}

	receipt, err := CreateDOTTransactionReceiptForTransaction(fundResult)
	if err != nil {
		return PaymentReceipt{}, err
	}

	uwa := UserPolkadotWalletAction{
		Action:             fmt.Sprintf("Funded crowdloan polkadot-%d for %f DOT", crowdloanId, value),
		UserUuid:           umw.UserUuid,
		CreatedAt:          time.Now(),
		PaymentReceiptUuid: receipt.Uuid,
	}

	return receipt, uwa.Save()
}

func (umw UserPolkadotWallet) WithdrawCrowdloan(crowdloanId uint64) (PaymentReceipt, error) {
	fundResult, err := apis.WithdrawCrowdloan(umw.PublicKey, apis.DOTWithdawCrowdloanMetadata{
		CrowdloanId: crowdloanId,
	})
	if err != nil {
		return PaymentReceipt{}, err
	}

	receipt, err := CreateDOTTransactionReceiptForTransaction(fundResult)
	if err != nil {
		return PaymentReceipt{}, err
	}

	uwa := UserPolkadotWalletAction{
		Action:             fmt.Sprintf("Withdrawed funds from crowdloan polkadot-%d", crowdloanId),
		UserUuid:           umw.UserUuid,
		CreatedAt:          time.Now(),
		PaymentReceiptUuid: receipt.Uuid,
	}

	return receipt, uwa.Save()
}

func (umw UserPolkadotWallet) PaybackCrowdloan(crowdloanId uint64) (PaymentReceipt, error) {
	fundResult, err := apis.PaybackCrowdloan(umw.PublicKey, apis.DOTPaybackCrowdloanMetadata{
		CrowdloanId: crowdloanId,
	})
	if err != nil {
		return PaymentReceipt{}, err
	}

	receipt, err := CreateDOTTransactionReceiptForTransaction(fundResult)
	if err != nil {
		return PaymentReceipt{}, err
	}

	uwa := UserPolkadotWalletAction{
		Action:             fmt.Sprintf("payback for crowdloan polkadot-%d", crowdloanId),
		UserUuid:           umw.UserUuid,
		CreatedAt:          time.Now(),
		PaymentReceiptUuid: receipt.Uuid,
	}

	return receipt, uwa.Save()
}

func (umw UserPolkadotWallet) PayoutCrowdloan(crowdloanId uint64) (PaymentReceipt, error) {
	fundResult, err := apis.CollectCrowdloanPayout(umw.PublicKey, apis.DOTPayoutCrowdloanMetadata{
		CrowdloanId: crowdloanId,
	})
	if err != nil {
		return PaymentReceipt{}, err
	}

	receipt, err := CreateDOTTransactionReceiptForTransaction(fundResult)
	if err != nil {
		return PaymentReceipt{}, err
	}

	uwa := UserPolkadotWalletAction{
		Action:             fmt.Sprintf("payout for crowdloan polkadot-%d", crowdloanId),
		UserUuid:           umw.UserUuid,
		CreatedAt:          time.Now(),
		PaymentReceiptUuid: receipt.Uuid,
	}

	return receipt, uwa.Save()
}

/*
	Queries
*/

func GetNextUserPolkadotWalletActionID() int {
	var userPolkadotWalletAction UserPolkadotWalletAction
	database.
		Order("ID desc").
		First(&userPolkadotWalletAction)

	return userPolkadotWalletAction.ID + 1
}

func GetNextUserPolkadotWalletbalanceID() int {
	var userPolkadotWalletAction UserPolkadotWalletBalance
	database.
		Order("ID desc").
		First(&userPolkadotWalletAction)

	return userPolkadotWalletAction.ID + 1
}

func GetAllUserPolkadotWallets() []UserPolkadotWallet {
	var items []UserPolkadotWallet
	database.Find(&items)
	return items
}

func GetAllUserPolkadotWalletBalances() []UserPolkadotWalletBalance {
	var items []UserPolkadotWalletBalance
	database.Find(&items)
	return items
}

func GetAllUserPolkadotWalletActions() []UserPolkadotWalletAction {
	var items []UserPolkadotWalletAction
	database.Find(&items)
	return items
}

func FindUserPolkadotWalletByPublicKey(publicKey string) (*UserPolkadotWallet, error) {
	var userPolkadotWallet UserPolkadotWallet

	err := database.
		Where(&UserPolkadotWallet{PublicKey: publicKey}).
		First(&userPolkadotWallet).
		Error
	if err != nil {
		return nil, err
	}

	return &userPolkadotWallet, err
}

func FindPolkadotWalletsForUser(userUuid string) []UserPolkadotWallet {
	var wallets []UserPolkadotWallet

	database.
		Where(&UserPolkadotWallet{UserUuid: userUuid}).
		Find(&wallets)

	return wallets
}

func FindUserPolkadotWalletActionsForUser(userUuid string) []UserPolkadotWalletAction {
	var actions []UserPolkadotWalletAction

	database.
		Where(&UserPolkadotWalletAction{UserUuid: userUuid}).
		Order("created_at ASC").
		Preload("PaymentReceipt").
		Find(&actions)

	for i := range actions {
		paymentReceipt := actions[i].PaymentReceipt
		if paymentReceipt.Type == "Polkadot" && paymentReceipt.SerializedData != "" {
			btcPaymentResult, err := paymentReceipt.BTCPaymentResult()
			if err != nil {
				continue
			}
			actions[i].PaymentReceipt.BTCPaymentResultItem = &btcPaymentResult
		}
	}

	return actions
}

// FindRecentPolkadotWallets returns wallets no older than 1 day
func FindRecentPolkadotWallets() []UserPolkadotWallet {
	var wallets []UserPolkadotWallet

	database.
		Where("created_at > ?", time.Now().AddDate(0, 0, -7)).
		Find(&wallets)

	return wallets
}

func FindAllPolkadotWallets() []UserPolkadotWallet {
	var wallets []UserPolkadotWallet
	database.Find(&wallets)
	return wallets
}

/*
	CRUD
*/

func CreatePolkadotWallet(user User) (*UserPolkadotWallet, error) {
	address, mnemonic, err := apis.GenerateDOTAddress("user_wallet")
	if err != nil {
		return nil, err
	}

	uw := UserPolkadotWallet{
		PublicKey: address,
		UserUuid:  user.Uuid,
		Mnemonic:  mnemonic,
	}

	uwa := UserPolkadotWalletAction{
		UserUuid:  uw.UserUuid,
		PublicKey: uw.PublicKey,
		Action:    "PolkadotWallet created",
		Amount:    float64(0.0),
	}
	uwa.Save()

	return &uw, uw.Save()
}

// Create views and other representatives
func setupUserPolkadotBalanceViews() {
	database.Exec("DROP VIEW IF EXISTS v_user_polkadot_wallet_balances CASCADE;")
	database.Exec(`
		CREATE VIEW v_user_polkadot_wallet_balances AS (
			select 
				sum(free_balance) as free_balance, 
				sum(reserved_balance) as reserved_balance, 
				sum(frozen_balance) as frozen_balance, 
				user_uuid, username  from (

					WITH UserPolkadotWalletBalancesUpdateTimes As (
					SELECT public_key, MAX(created_at) max_timestamp
					FROM user_polkadot_wallet_balances
					GROUP BY public_key
					)
					select 
						uwb.created_at, uwb.public_key, 
						uwb.free_balance, uwb.reserved_balance, uwb.frozen_balance, 
						uw.user_uuid, u.username
					from 
						user_polkadot_wallet_balances uwb 
					join 
						user_polkadot_wallets uw on uw.public_key=uwb.public_key 
					join
						users u on u.uuid=uw.user_uuid
					inner join 
						UserPolkadotWalletBalancesUpdateTimes t on t.public_key=uwb.public_key and uwb.created_at = t.max_timestamp
					) uwb
					group by username, user_uuid
					order by free_balance desc
	);`)
}
