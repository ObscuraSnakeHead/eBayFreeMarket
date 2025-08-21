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

type UserMoneroWallet struct {
	PublicKey string `json:"public_key" gorm:"primary_key"`
	UserUuid  string `json:"-" gorm:"index"`
	IsLocked  bool   `json:"-"`

	User User `json:"-"`

	CreatedAt *time.Time `json:"created_at,omitempty"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}

type UserMoneroWalletBalance struct {
	ID              int     `json:"id" gorm:"primary_key"`
	PublicKey       string  `json:"public_key" gorm:"index"`
	Balance         float64 `json:"balance"`
	UnlockedBalance float64 `json:"unlocked_balance"`
	Type            string  `json:"type"`

	CreatedAt time.Time `json:"created_at"`

	UserMoneroWallet UserMoneroWallet `gorm:"AssociationForeignKey:PublicKey" json:"-"`
}

type UserMoneroWallets []UserMoneroWallet

type UserMoneroWalletAction struct {
	ID                 int    `json:"id" gorm:"primary_key"`
	UserUuid           string `json:"-" gorm:"user_uuid"`
	PaymentReceiptUuid string `json:"-" gorm:"payment_receipt_uuid"`
	PublicKey          string `json:"-" gorm:"index"`

	Action string  `json:"action"`
	Amount float64 `json:"amount"`

	CreatedAt time.Time `json:"created_at"`

	User             User             `json:"-"`
	UserMoneroWallet UserMoneroWallet `gorm:"AssociationForeignKey:PublicKey" json:"-"`
	PaymentReceipt   PaymentReceipt   `json:"payment_receipt"`
}

/*
	Model Methods
*/

func (w UserMoneroWallet) UpdateBalance(force bool) (UserMoneroWalletBalance, error) {
	d, _ := time.ParseDuration("1m")
	now := time.Now()

	if w.UpdatedAt != nil && !w.UpdatedAt.Add(d).Before(now) && !force {
		return UserMoneroWalletBalance{}, errors.New("MoneroWallet was updated recently. Please wait.")
	}

	w.UpdatedAt = &now

	err := w.Save()
	if err != nil {
		return UserMoneroWalletBalance{}, err
	}

	balance, err := apis.GetAmountOnXMRAddress(w.PublicKey)
	if err != nil {
		return UserMoneroWalletBalance{}, err
	}

	currentBalance := w.CurrentBalance()
	if currentBalance.Balance != balance.Balance || currentBalance.UnlockedBalance != balance.UnlockedBalance {

		uwb := UserMoneroWalletBalance{
			PublicKey:       w.PublicKey,
			Balance:         balance.Balance,
			UnlockedBalance: balance.UnlockedBalance,
			CreatedAt:       now,
			Type:            "blockchain",
		}
		uwb.Save()

		var (
			action string
			amount float64
		)

		if currentBalance.Balance != balance.Balance {
			action = "Balance updated"
			amount = balance.Balance
		} else {
			action = "Unlocked balance updated"
			amount = balance.UnlockedBalance
		}

		uwa := UserMoneroWalletAction{
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

func (w UserMoneroWallet) CurrentBalance() UserMoneroWalletBalance {
	var uwb UserMoneroWalletBalance

	database.
		Where(&UserMoneroWalletBalance{PublicKey: w.PublicKey}).
		Order("created_at DESC").
		First(&uwb)

	return uwb
}

/*
	Database methods
*/

// UserMoneroWallet

func (w UserMoneroWallet) Validate() error {
	if w.PublicKey == "" {
		return errors.New("Wrong wallet")
	}
	return nil
}

func (w UserMoneroWallet) Save() error {
	err := w.Validate()
	if err != nil {
		return err
	}
	return w.SaveToDatabase()
}

func (w UserMoneroWallet) Remove() error {
	return database.Delete(w).Error
}

func (w UserMoneroWallet) SaveToDatabase() error {
	if existing, _ := FindUserMoneroWalletByPublicKey(w.PublicKey); existing == nil {
		return database.Create(&w).Error
	}
	return database.Save(&w).Error
}

// UserMoneroWalletBalance

func (w UserMoneroWalletBalance) Validate() error {
	if w.PublicKey == "" {
		return errors.New("Wrong wallet")
	}
	return nil
}

func (w UserMoneroWalletBalance) Save() error {
	err := w.Validate()
	if err != nil {
		return err
	}
	return w.SaveToDatabase()
}

func (w UserMoneroWalletBalance) Remove() error {
	return database.Delete(w).Error
}

func (w UserMoneroWalletBalance) SaveToDatabase() error {
	if w.ID != 0 {
		return database.Save(&w).Error
	}
	w.ID = GetNextUserMoneroWalletbalanceID()
	return database.Create(&w).Error
}

func (w UserMoneroWalletAction) Validate() error {
	if w.UserUuid == "" {
		return errors.New("Wrong UserUuid")
	}
	return nil
}

func (w UserMoneroWalletAction) Save() error {
	err := w.Validate()
	if err != nil {
		return err
	}
	return w.SaveToDatabase()
}

func (w UserMoneroWalletAction) Remove() error {
	return database.Delete(w).Error
}

func (w UserMoneroWalletAction) SaveToDatabase() error {
	if w.ID != 0 {
		return database.Save(&w).Error
	}
	w.ID = GetNextUserMoneroWalletActionID()
	return database.Create(&w).Error
}

/*
	Model Methods
*/

func (umws UserMoneroWallets) Balance() apis.XMRWalletBalance {
	var balance apis.XMRWalletBalance

	for _, umw := range umws {
		balance.Balance += umw.CurrentBalance().Balance
		balance.UnlockedBalance += umw.CurrentBalance().UnlockedBalance
	}

	return balance
}

func (umw UserMoneroWallet) Balance() apis.XMRWalletBalance {
	var balance apis.XMRWalletBalance

	balance.Balance = umw.CurrentBalance().Balance
	balance.UnlockedBalance = umw.CurrentBalance().UnlockedBalance

	return balance
}

func (w UserMoneroWallet) prepareTransaction(address string, amount float64) (apis.XMRPayment, error) {
	if amount > w.Balance().Balance {
		return apis.XMRPayment{}, errors.New("Amount is greater than balance")
	}

	if !moneroRegexp.MatchString(address) {
		return apis.XMRPayment{}, errors.New("Wrong XMR address")
	}

	if w.IsLocked {
		return apis.XMRPayment{}, errors.New("Wallet is currently locked")
	}
	cb := w.CurrentBalance()
	if cb.UnlockedBalance == float64(0.0) {
		return apis.XMRPayment{}, errors.New("Unlocked balance is 0. Please wait until funds are unlocked.")
	}

	return apis.XMRPayment{Address: address, Amount: amount}, nil
}

func (umw UserMoneroWallet) Send(address string, amount float64) (PaymentReceipt, error) {
	payment, err := umw.prepareTransaction(address, amount)
	if err != nil {
		return PaymentReceipt{}, err
	}

	transactionResult, err := apis.SendXMR(address, []apis.XMRPayment{payment})
	if err != nil {
		return PaymentReceipt{}, err
	}

	receipt, err := CreateXMRPaymentReceipt(transactionResult)
	if err != nil {
		return PaymentReceipt{}, err
	}

	uwa := UserMoneroWalletAction{
		Action:             fmt.Sprintf("Sent %f XMR to %s", amount, address),
		UserUuid:           umw.UserUuid,
		CreatedAt:          time.Now(),
		PaymentReceiptUuid: receipt.Uuid,
	}

	return receipt, uwa.Save()
}

/*
	Queries
*/

func GetNextUserMoneroWalletActionID() int {
	var userMoneroWalletAction UserMoneroWalletAction
	database.
		Order("ID desc").
		First(&userMoneroWalletAction)

	return userMoneroWalletAction.ID + 1
}

func GetNextUserMoneroWalletbalanceID() int {
	var userMoneroWalletAction UserMoneroWalletBalance
	database.
		Order("ID desc").
		First(&userMoneroWalletAction)

	return userMoneroWalletAction.ID + 1
}

func GetAllUserMoneroWallets() []UserMoneroWallet {
	var items []UserMoneroWallet
	database.Find(&items)
	return items
}

func GetAllUserMoneroWalletBalances() []UserMoneroWalletBalance {
	var items []UserMoneroWalletBalance
	database.Find(&items)
	return items
}

func GetAllUserMoneroWalletActions() []UserMoneroWalletAction {
	var items []UserMoneroWalletAction
	database.Find(&items)
	return items
}

func FindUserMoneroWalletByPublicKey(publicKey string) (*UserMoneroWallet, error) {
	var userMoneroWallet UserMoneroWallet

	err := database.
		Where(&UserMoneroWallet{PublicKey: publicKey}).
		First(&userMoneroWallet).
		Error
	if err != nil {
		return nil, err
	}

	return &userMoneroWallet, err
}

func FindMoneroWalletsForUser(userUuid string) []UserMoneroWallet {
	var wallets []UserMoneroWallet

	database.
		Where(&UserMoneroWallet{UserUuid: userUuid}).
		Find(&wallets)

	return wallets
}

func FindUserMoneroWalletActionsForUser(userUuid string) []UserMoneroWalletAction {
	var actions []UserMoneroWalletAction

	database.
		Where(&UserMoneroWalletAction{UserUuid: userUuid}).
		Order("created_at ASC").
		Preload("PaymentReceipt").
		Find(&actions)

	for i := range actions {
		paymentReceipt := actions[i].PaymentReceipt
		if paymentReceipt.Type == "Monero" && paymentReceipt.SerializedData != "" {
			btcPaymentResult, err := paymentReceipt.BTCPaymentResult()
			if err != nil {
				continue
			}
			actions[i].PaymentReceipt.BTCPaymentResultItem = &btcPaymentResult
		}
	}

	return actions
}

// FindRecentMoneroWallets returns wallets no older than 1 day
func FindRecentMoneroWallets() []UserMoneroWallet {
	var wallets []UserMoneroWallet

	database.
		Where("created_at > ?", time.Now().AddDate(0, 0, -7)).
		Find(&wallets)

	return wallets
}

func FindAllMoneroWallets() []UserMoneroWallet {
	var wallets []UserMoneroWallet
	database.Find(&wallets)
	return wallets
}

/*
	CRUD
*/

func CreateMoneroWallet(user User) (*UserMoneroWallet, error) {
	address, err := apis.GenerateXMRAddress("user_wallet")
	if err != nil {
		return nil, err
	}

	uw := UserMoneroWallet{
		PublicKey: address,
		UserUuid:  user.Uuid,
	}

	uwa := UserMoneroWalletAction{
		UserUuid:  uw.UserUuid,
		PublicKey: uw.PublicKey,
		Action:    "MoneroWallet created",
		Amount:    float64(0.0),
	}
	uwa.Save()

	return &uw, uw.Save()
}

// Create views and other representatives
func setupUserMoneroBalanceViews() {
	database.Exec("DROP VIEW IF EXISTS v_user_monero_wallet_balances CASCADE;")
	database.Exec(`
		CREATE VIEW v_user_monero_wallet_balances AS (
			select sum(balance) as balance, sum(unlocked_balance) as unlocked_balance, user_uuid, username  from (

			WITH UserMoneroWalletBalancesUpdateTimes As (
			   SELECT public_key, MAX(created_at) max_timestamp
			   FROM user_monero_wallet_balances
			   GROUP BY public_key
			)
			select 
				uwb.created_at, uwb.public_key, uwb.balance, uwb.unlocked_balance, uw.user_uuid, u.username
			from 
				user_monero_wallet_balances uwb 
			join 
				user_monero_wallets uw on uw.public_key=uwb.public_key 
			join
				users u on u.uuid=uw.user_uuid
			inner join 
				UserMoneroWalletBalancesUpdateTimes t on t.public_key=uwb.public_key and uwb.created_at = t.max_timestamp
			) uwb
			group by username, user_uuid
			order by balance desc
	);`)
}
