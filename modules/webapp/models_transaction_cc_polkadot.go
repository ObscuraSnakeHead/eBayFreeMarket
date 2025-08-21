package webapp

import (
	"errors"
	"fmt"

	"ucoi4tzlito52agmquc6oopn3zpmr6djz5vvfabtgrhyc6hufpzjtnad.onion/Tochka/tochka-free-market/modules/apis"
)

type PolkadotTransaction struct {
	Uuid   string  `json:"uuid" gorm:"primary_key"`
	Id     uint64  `json:"id" gorm:"index"`
	Amount float64 `json:"amount"`
}

/*
	Database Methods
*/

func (t PolkadotTransaction) Save() error {
	if existing, _ := FindPolkadotTransactionById(t.Id); existing == nil {
		return database.Create(&t).Error
	}
	return database.Save(&t).Error
}

func FindPolkadotTransactionById(id uint64) (*PolkadotTransaction, error) {
	var item PolkadotTransaction
	err := database.
		Where(&PolkadotTransaction{Id: id}).
		First(&item).
		Error
	if err != nil {
		return nil, err
	}
	return &item, err
}

// Contract Methods

func CreatePolkadotTransaction(
	itemPackage Package,
	buyer User,
	tp string,
	quantity int,
	shippingPrice float64,
) (PolkadotTransaction, error) {
	buyerPolkadotWallets := buyer.FindUserPolkadotWallets()
	if len(buyerPolkadotWallets) == 0 {
		return PolkadotTransaction{}, errors.New("Buyer doesn't have Polkadot Wallet")
	}
	buyerPolkadotWallet := buyerPolkadotWallets[0]

	sellerPolkadotWallets := itemPackage.Store.AdminUser().FindUserPolkadotWallets()
	if len(sellerPolkadotWallets) == 0 {
		return PolkadotTransaction{}, errors.New("Buyer doesn't have Polkadot Wallet")
	}
	sellerPolkadotWallet := sellerPolkadotWallets[0]

	dotValue := itemPackage.GetPrice("ASTR")*float64(quantity) + shippingPrice

	transactionResult, err := apis.MintEscrow(buyerPolkadotWallet.PublicKey, apis.DOTEscrowMintMetadata{
		SellerAddress: sellerPolkadotWallet.PublicKey,
		Cid:           "",
		Value:         dotValue,
	})

	polkadotTransaction := PolkadotTransaction{
		Uuid:   fmt.Sprintf("polkadot-escrow-%d", transactionResult.Id),
		Id:     transactionResult.Id,
		Amount: dotValue,
	}
	polkadotTransaction.Save()

	receipt, err := CreateDOTEscrowReceipt([]apis.DOTEscrowMintResult{transactionResult})
	if err != nil {
		return polkadotTransaction, err
	}

	uwa := UserPolkadotWalletAction{
		UserUuid:           buyer.Uuid,
		PublicKey:          buyerPolkadotWallet.PublicKey,
		PaymentReceiptUuid: receipt.Uuid,
		Action:             "User minted escrow",
		Amount:             dotValue,
	}

	return polkadotTransaction, uwa.Save()
}

/*
	Financial Methods
*/

func (bt PolkadotTransaction) UpdateTransactionStatus(t Transaction) error {
	escrowInfo, err := apis.EscrowInfo(bt.Id)
	if err != nil {
		return err
	}

	if t.CurrentPaymentStatus() != escrowInfo.EscrowStatus.Status {
		err = t.SetTransactionStatus(escrowInfo.EscrowStatus.Status, float64(escrowInfo.Value), "Escrow record updated", "", nil)
		if err != nil {
			return err
		}
	}

	if t.CurrentShippingStatus() != escrowInfo.ShippingStatus.Status {
		err = t.SetShippingStatus(escrowInfo.ShippingStatus.Status, "Updated from blockchain", "")
		if err != nil {
			return err
		}
	}

	return nil
}

func (bt PolkadotTransaction) PartialRefund(t Transaction, comment, userUuid string, refundPercent float64) error {
	// TODO: Implement
	return nil
}

func (bt PolkadotTransaction) Release(t Transaction, comment, userUuid string) error {
	transactionResult, err := apis.ReleaseEscrow(bt.Id)
	if err != nil {
		return err
	}

	receipt, err := CreateDOTTransactionReceiptForTransaction(transactionResult)
	if err != nil {
		return err
	}

	return t.SetTransactionStatus(
		"RELEASED",
		t.CurrentAmountPaid(),
		comment,
		userUuid,
		&receipt,
	)
}

func (bt PolkadotTransaction) Cancel(t Transaction, comment, userUuid string) error {
	transactionResult, err := apis.CancelEscrow(bt.Id)
	if err != nil {
		return err
	}

	receipt, err := CreateDOTTransactionReceiptForTransaction(transactionResult)
	if err != nil {
		return err
	}

	return t.SetTransactionStatus(
		"CANCELLED",
		t.CurrentAmountPaid(),
		comment,
		userUuid,
		&receipt,
	)
}
