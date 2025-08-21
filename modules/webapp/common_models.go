package webapp

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var database *gorm.DB

func SyncModels() {
	database.AutoMigrate(
		&Advertising{},
		&APISession{},
		&BitcoinTransaction{},
		&City{},
		&CityMetroStation{},
		&Country{},
		&Crowdloan{},
		&CrowdloanStatus{},
		&Deposit{},
		&DepositHistory{},
		&Dispute{},
		&DisputeClaim{},
		&EthereumTransaction{},
		&Item{},
		&ItemCategory{},
		&Message{},
		&MoneroTransaction{},
		&Package{},
		&PackagePrice{},
		&PaymentReceipt{},
		&PolkadotTransaction{},
		&RatingReview{},
		&ReferralPayment{},
		&Reservation{},
		&ShippingOption{},
		&ShippingStatus{},
		&Store{},
		&StoreUser{},
		&StoreWarning{},
		&SupportTicket{},
		&SupportTicketStatus{},
		&ThreadPerusalStatus{},
		&Transaction{},
		&TransactionStatus{},
		&User{},
		&UserBitcoinWallet{},
		&UserBitcoinWalletAction{},
		&UserBitcoinWalletBalance{},
		&UserEthereumWallet{},
		&UserEthereumWalletAction{},
		&UserEthereumWalletBalance{},
		&UserMoneroWallet{},
		&UserMoneroWalletAction{},
		&UserMoneroWalletBalance{},
		&UserPolkadotWallet{},
		&UserPolkadotWalletAction{},
		&UserPolkadotWalletBalance{},
		&UserSettingsHistory{},
	)
}

func SyncDatabaseViews() {
	// drop all views and triggers

	database.Exec(`
		CREATE OR REPLACE FUNCTION strip_all_triggers() RETURNS text AS $$ DECLARE
	    triggNameRecord RECORD;
	    triggTableRecord RECORD;
	BEGIN
	    FOR triggNameRecord IN select distinct(trigger_name) from information_schema.triggers where trigger_schema = 'public' LOOP
	        FOR triggTableRecord IN SELECT distinct(event_object_table) from information_schema.triggers where trigger_name = triggNameRecord.trigger_name LOOP
	            RAISE NOTICE 'Dropping trigger: % on table: %', triggNameRecord.trigger_name, triggTableRecord.event_object_table;
	            EXECUTE 'DROP TRIGGER ' || triggNameRecord.trigger_name || ' ON ' || triggTableRecord.event_object_table || ';';
	        END LOOP;
	    END LOOP;

	    RETURN 'done';
	END;
	$$ LANGUAGE plpgsql SECURITY DEFINER;
	`)

	database.Exec(`
		select strip_all_triggers();
	`)

	database.Exec(`
	SELECT 
	'DROP VIEW ' || table_name || ';'
	FROM information_schema.views
	WHERE table_schema NOT IN ('pg_catalog', 'information_schema')
	AND table_name !~ '^pg_';
	`)

	// wallets & balances
	setupUserBitcoinBalanceViews()
	setupUserEthereumBalanceViews()
	setupUserMoneroBalanceViews()
	setupUserPolkadotBalanceViews()

	// messageboard & messages
	setupThreadsViews()
	setupPrivateThreadsFunctions()
	setupVendorVerificationThreadsFunctions()
	// setupMessageboardCategoriesViews()

	// transcations
	setupTransactionStatusesView()

	// users
	setupUserViews()
	setupVendorTxStatsViews()
	setupItemTxStatsViews()

	// items & packages, categories
	setupCategoriesViews()
	setupPackagesView()
	setupSerpItemsView()

	// tickets
	setupSupportTicketViews()

	// advertisings
	setupAdvertisingViews()
}

func init() {
	var err error

	database, err = gorm.Open("postgres", MARKETPLACE_SETTINGS.PostgresConnectionString)
	if err != nil {
		panic(err)
	}
	database.DB().SetMaxOpenConns(40)
	database.DB().Ping()
}
