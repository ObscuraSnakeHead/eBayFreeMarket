package webapp

import (
	"ucoi4tzlito52agmquc6oopn3zpmr6djz5vvfabtgrhyc6hufpzjtnad.onion/Tochka/tochka-free-market/modules/apis"
)

type MoneroTransaction struct {
	Uuid   string  `json:"uuid" gorm:"primary_key"`
	Amount float64 `json:"amount"`
}

/*
	Database MXMRods
*/

func (t MoneroTransaction) Save() error {
	if existing, _ := FindMoneroTransactionByUuid(t.Uuid); existing == nil {
		return database.Create(&t).Error
	}
	return database.Save(&t).Error
}

func FindMoneroTransactionByUuid(uuid string) (*MoneroTransaction, error) {
	var item MoneroTransaction
	err := database.
		Where(&MoneroTransaction{Uuid: uuid}).
		First(&item).
		Error
	if err != nil {
		return nil, err
	}
	return &item, err
}

/*
	Financial MXMRods
*/

func (bt MoneroTransaction) UpdateTransactionStatus(t Transaction) error {
	newAmount, err := apis.GetAmountOnXMRAddress(t.Uuid)
	if err != nil {
		return err
	}
	if t.CurrentAmountPaid() == newAmount.UnlockedBalance {
		return nil
	}
	if t.IsPending() {
		if (bt.Amount - newAmount.UnlockedBalance) <= bt.Amount {
			return t.SetTransactionStatus(
				"COMPLETED",
				newAmount.Balance,
				"Transaction funded",
				"",
				nil,
			)
		}
		return t.SetTransactionStatus(
			"PENDING",
			newAmount.Balance,
			"Transaction amount updated",
			"",
			nil,
		)
	}
	return t.SetTransactionStatus(
		t.CurrentPaymentStatus(),
		newAmount.Balance,
		"Transaction amount updated",
		"",
		nil,
	)
}

func (bt MoneroTransaction) PartialRefund(t Transaction, comment, userUuid string, refundPercent float64) error {
	var (
		addressFrom = t.Uuid
		buyer       = t.Buyer
		store       = t.Store
		err         error
	)

	resolver, err := FindUserByUuid(userUuid)
	if err != nil {
		return err
	}

	storeAddress := store.MoneroAddress()
	buyerAddress := buyer.MoneroAddress()
	resolverAddress := resolver.MoneroAddress()

	payments := []apis.XMRPayment{
		{
			Address: storeAddress,
			Percent: (1. - refundPercent) - 0.05,
		},
		{
			Address: buyerAddress,
			Percent: refundPercent - 0.05,
		},
		{
			Address: resolverAddress,
			Percent: 0.1,
		},
	}

	results, err := apis.SendXMR(addressFrom, payments)
	if err != nil {
		return err
	}

	receipt, err := CreateXMRPaymentReceipt(results)
	if err != nil {
		return err
	}

	if len(t.Status) > 0 && t.Status[(len(t.Status)-1)].Status != "CANCELLED" {
		t.SetTransactionStatus(
			"CANCELLED",
			t.CurrentAmountPaid(),
			comment,
			userUuid,
			&receipt,
		)
	}

	return nil
}

func (bt MoneroTransaction) Release(t Transaction, comment, userUuid string) error {
	var (
		addressFrom           = t.Uuid
		buyer                 = t.Buyer
		buyerReferralPayment  *ReferralPayment
		commission            = t.CommissionPercent()
		payments              = []apis.XMRPayment{}
		usdRate               = GetCurrencyRate("XMR", "USD")
		store                 = t.Store
		vendorReferralPayment *ReferralPayment
		err                   error
	)

	// Vendor address
	addressTo := store.MoneroAddress()

	payments = []apis.XMRPayment{
		{
			Address: addressTo,
			Percent: 1. - commission,
		},
	}

	var buyerInviter *User
	if t.MobileAppTransaction {
		buyerInviter, _ = FindUserByUsername(MARKETPLACE_SETTINGS.AndroidDeveloperUsername)
	} else {
		buyerInviter = buyer.Iniviter()
	}

	if buyerInviter != nil {
		inviterPercent := MARKETPLACE_SETTINGS.FreeAccountReferralPercent
		if buyerInviter.IsGoldAccount {
			inviterPercent = MARKETPLACE_SETTINGS.GoldAccountCommission
		}
		if buyerInviter.IsSilverAccount {
			inviterPercent = MARKETPLACE_SETTINGS.SilverAccountCommission
		}
		if buyerInviter.IsBronzeAccount {
			inviterPercent = MARKETPLACE_SETTINGS.BronzeAccountCommission
		}
		if t.MobileAppTransaction {
			inviterPercent = MARKETPLACE_SETTINGS.AndroidDeveloperCommission
		}

		inviterWallet := buyerInviter.Monero
		if inviterWallet == "" {
			inviterWallet = buyerInviter.FindRecentMoneroWallet().PublicKey
		}

		payments = append(payments, apis.XMRPayment{
			Address: inviterWallet,
			Percent: commission * inviterPercent,
		})

		buyerReferralPayment = &ReferralPayment{
			TransactionUuid:    t.Uuid,
			ReferralPercent:    inviterPercent,
			ReferralPaymentBTC: commission * inviterPercent * bt.Amount,
			ReferralPaymentUSD: commission * inviterPercent * bt.Amount * usdRate,
			UserUuid:           buyer.InviterUuid,
			IsBuyerReferral:    true,
		}
	}

	commissionPercent := 1.0
	for _, p := range payments {
		commissionPercent -= p.Percent
	}

	payments = append(payments, apis.XMRPayment{
		Address: MARKETPLACE_SETTINGS.MoneroCommissionWallet,
		Percent: commissionPercent,
	})

	results, err := apis.SendXMR(addressFrom, payments)
	if err != nil {
		return err
	}

	receipt, err := CreateXMRPaymentReceipt(results)
	if err != nil {
		return err
	}

	if len(t.Status) > 0 && t.Status[(len(t.Status)-1)].Status != "RELEASED" {
		t.SetTransactionStatus(
			"RELEASED",
			t.CurrentAmountPaid(),
			comment,
			userUuid,
			&receipt,
		)
		if buyerReferralPayment != nil {
			buyerReferralPayment.Save()
		}
		if vendorReferralPayment != nil {
			vendorReferralPayment.Save()
		}
	}

	return nil
}

func (bt MoneroTransaction) Cancel(t Transaction, comment, userUuid string) error {
	buyer, err := FindUserByUuid(t.BuyerUuid)
	if err != nil {
		return err
	}

	var (
		addressFrom = t.Uuid
		buyerWallet = buyer.FindRecentMoneroWallet()
		addressTo   = buyerWallet.PublicKey
	)

	payments := []apis.XMRPayment{
		{Address: addressTo, Percent: 1.0},
	}

	results, err := apis.SendXMR(addressFrom, payments)
	if err != nil {
		return err
	}

	receipt, err := CreateXMRPaymentReceipt(results)
	if err != nil {
		return err
	}

	if len(t.Status) > 0 && t.Status[(len(t.Status)-1)].Status != "CANCELLED" {
		t.SetTransactionStatus(
			"CANCELLED",
			t.CurrentAmountPaid(),
			comment,
			userUuid,
			&receipt,
		)
	}

	return nil
}

func CreateMoneroTransaction(
	itemPackage Package,
	buyer User,
	tp string,
	quantity int,
	shippingPrice float64,
) (MoneroTransaction, error) {
	wallet, err := apis.GenerateXMRAddress("escrow")
	if err != nil {
		return MoneroTransaction{}, err
	}

	MoneroTransaction := MoneroTransaction{
		Uuid:   wallet,
		Amount: itemPackage.GetPrice("XMR")*float64(quantity) + shippingPrice,
	}

	return MoneroTransaction, MoneroTransaction.Save()
}

/*
	Tx Stats
*/

func GetMoneroTxStatsForVendor(uuid string) TxStats {
	var stats TxStats

	database.
		Table("v_current_monero_transaction_statuses").
		Select("count(*) as tx_number, sum(amount) as tx_volume").
		Where("store_uuid = ?", uuid).
		Where("current_status NOT IN ('CANCELLED', 'FAILED', 'PENDING')").
		First(&stats)

	return stats
}

func GetMoneroTxStatsForBuyer(uuid string) TxStats {
	var stats TxStats

	database.
		Table("v_current_monero_transaction_statuses").
		Select("count(*) as tx_count, sum(amount) as tx_volume").
		Where("buyer_uuid = ?", uuid).
		Where("current_status NOT IN ('CANCELLED', 'FAILED', 'PENDING')").
		First(&stats)

	return stats
}
