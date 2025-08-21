package webapp

import (
	"github.com/helloyi/go-waitgroup"
	"github.com/robfig/cron"

	"ucoi4tzlito52agmquc6oopn3zpmr6djz5vvfabtgrhyc6hufpzjtnad.onion/Tochka/tochka-free-market/modules/util"
)

// TaskCleanInactiveReservations is a cron job that runs every 5 minutes.
// Deletes inactive reservations.
func TaskCleanInactiveReservations() {
	inactiveReservations := FindInactiveReservations()
	util.Log.Debug("[Task] [TaskCleanInactiveReservations] # of items:", len(inactiveReservations))
	for _, r := range inactiveReservations {
		transaction := r.Transaction
		transaction.RefreshTransactionStatus()
		if transaction.IsFailed() && transaction.IsCancelled() {
			transaction.SetTransactionStatus(
				"FAILED",
				transaction.CurrentAmountPaid(),
				"Transaction failed because no coins were transferred",
				"",
				nil,
			)
		}
		r.Remove()
	}
}

// TaskUpdatePendingTransactions is a cron job that runs every 5 minutes.
// Checks for balance and updates transaction status.
func TaskUpdatePendingTransactions() {
	it := FindPendingTransactions()
	util.Log.Debug("[Task] [TaskUpdatePendingTransactions] # of items:", len(it))
	for _, t := range it {
		t.RefreshTransactionStatus()
		t.UpdateTransactionStatus()
	}
}

// TaskFailOldPendingTransactions is a cron job that runs every 5 minutes.
// Updates status of old pending payments
func TaskFailOldPendingTransactions() {
	it := FindOldPendingTransactions()
	util.Log.Debug("[Task] [TaskFailOldPendingTransactions] # of items:", len(it))
	for _, t := range it {
		t.Fail("Escrow failed automatically", "")
	}
}

// TaskReleaseConfirmedTransactions is a cron job that runs every 5 minutes.
// Releases funds to seller of confirmed transaction.
func TaskReleaseOldCompletedTransactions() {
	it := FindOldCompletedTransactions()
	util.Log.Debug("[Task] [TaskReleaseOldCompletedTransactions] # of items:", len(it))
	for _, t := range it {
		t.Release("Escrow released automatically", "")
	}
}

func TaskUpdateBalancesOrRecentlyReleasedAndCancelledTransactions() {
	wg := waitgroup.Create(16)
	ts := FindRecentlyCancelledAndReleasedTransactions()
	util.Log.Debug("[Task] [TaskUpdateBalancesOrRecentlyReleasedAndCancelledTransactions] # of items:", len(ts))

	for i := range ts {
		t := ts[i]
		wg.BlockAdd()
		go func(t *Transaction) {
			defer wg.Done()
			t.RefreshTransactionStatus()
			t.UpdateTransactionStatus()
		}(&t)
	}

	wg.Wait()
}

func TaskFinalizeReleasedAndCancelledTransactionsWithNonZeroAmount() {
	it := FindReleasedAndCancelledTransactionsWithNonZeroAmount()
	wg := waitgroup.Create(16)

	util.Log.Debug("[Task] [TaskFinalizeReleasedAndCancelledTransactionsWithNonZeroAmount] # of items:", len(it))
	for i := range it {
		t := it[i]
		wg.BlockAdd()
		go func(t *Transaction) {
			defer wg.Done()
			t.RefreshTransactionStatus()
			if t.CurrentPaymentStatus() == "CANCELLED" {
				t.Cancel("Tx cancelled automatically", "")
			}
			if t.CurrentPaymentStatus() == "RELEASED" {
				t.Release("Tx released automatically", "")
			}
			t.UpdateTransactionStatus()
		}(&t)
	}

	wg.Wait()
}

func TaskFreezeStuckCompletedTransactions() {
	it := FindStuckCompletedTransactions()
	util.Log.Debug("[Task] [TaskFreezeStuckCompletedTransactions] # of items:", len(it))
	for _, t := range it {
		t.RefreshTransactionStatus()
		t.Freeze("Escrow frozen because transaction has failed to release", "")
	}
}

func TaskFreezePre2024Transactions() {
	it := FindPre2024CompletedAndPendingTransactions()
	util.Log.Debug("[Task] [TaskFreezePre2024Transactions] # of items:", len(it))
	for _, t := range it {
		t.RefreshTransactionStatus()
		t.Freeze("Escrow frozen because of service suspension in 2019", "")
	}
}

func CancelCompletedAndNotDispatchedTransactions() {
	it := FindCompletedAndNotDispatchedTransactions()
	util.Log.Debug("[Task] [CancelCompletedAndNotDispatchedTransactions] # of items:", len(it))
	for _, t := range it {
		t.RefreshTransactionStatus()
		t.Freeze("Transaction has been cancelled because it hasn't been marked as dispatched within 3 days", "")
	}
}

func StartTransactionsCron() {
	c := cron.New()
	c.AddFunc("@every 10m", TaskCleanInactiveReservations)
	c.AddFunc("@every 2m", TaskUpdatePendingTransactions)
	c.AddFunc("@every 10m", TaskFailOldPendingTransactions)
	// c.AddFunc("@every 10m", TaskReleaseOldCompletedTransactions)
	c.AddFunc("@every 10m", TaskUpdateBalancesOrRecentlyReleasedAndCancelledTransactions)
	// c.AddFunc("@hourly", TaskFinalizeReleasedAndCancelledTransactionsWithNonZeroAmount)
	c.AddFunc("@hourly", TaskFreezeStuckCompletedTransactions)
	c.AddFunc("@hourly", CancelCompletedAndNotDispatchedTransactions)

	// Materialized Views
	c.AddFunc("@every 1m", RefreshShippingStatusMaterializedView)
	c.AddFunc("@every 1m", RefreshThreadCountMaterializedView)

	c.Start()
}
