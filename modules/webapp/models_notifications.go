package webapp

import (
	"math/rand"
	"time"

	"ucoi4tzlito52agmquc6oopn3zpmr6djz5vvfabtgrhyc6hufpzjtnad.onion/Tochka/tochka-free-market/modules/util"
)

/*
	Models
*/

type Notification struct {
	Uuid string `json:"-" gorm:"primary_key"`

	PrivateMessages  int `json:"private_messages"`
	Purchases        int `json:"transactions"`
	PurchaseMessages int `json:"transaction_messages"`
	Disputes         int `json:"disputes"`
	SupportMessages  int `json:"support_messages"`

	NewPrivateMessages  int `json:"new_private_messages"`
	NewPurchases        int `json:"new_transactions"`
	NewPurchaseMessages int `json:"new_transaction_messages"`
	NewDisputes         int `json:"new_disputes"`
	NewSupportMessages  int `json:"new_support_messages"`

	// TimeStamps
	CreatedAt *time.Time `json:"-" gorm:"index"`
	UpdatedAt *time.Time `json:"-" gorm:"index"`
	DeletedAt *time.Time `json:"-" gorm:"index"`
}

func NewRandomNotification() Notification {
	return Notification{
		Uuid:                util.GenerateUuid(),
		NewPrivateMessages:  rand.Intn(1000),
		NewPurchases:        rand.Intn(1000),
		NewPurchaseMessages: rand.Intn(1000),
	}
}

func GetCurrentNotification(user User) Notification {
	store := user.Store()
	storeUuid := ""
	if store != nil {
		storeUuid = store.Uuid
	}
	numberOfTx := CountNumberOfTransactions(user.Uuid, storeUuid)

	return Notification{
		PrivateMessages:     CountPrivateMessages(user),
		Disputes:            CountDisputesForUserUuid(user.Uuid, ""),
		SupportMessages:     CountSupportTicketsForUser(user),
		Purchases:           numberOfTx,
		PurchaseMessages:    100,
		NewPrivateMessages:  CountUndreadPrivateMessages(user),
		NewSupportMessages:  100,
		NewPurchases:        100,
		NewPurchaseMessages: 100,
		NewDisputes:         100,
	}
}
