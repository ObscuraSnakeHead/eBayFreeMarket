package webapp

import (
	"github.com/gocraft/web"

	"ucoi4tzlito52agmquc6oopn3zpmr6djz5vvfabtgrhyc6hufpzjtnad.onion/Tochka/tochka-free-market/modules/settings"
)

func ConfigureRouter(router *web.Router) *web.Router {
	// ----------
	// Middleware
	// ----------

	router.Middleware((*Context).PerformanceMeasureMiddleware)
	router.Middleware((*Context).LoggerMiddleware)

	if settings.GetSettings().BotCheck {
		router.Middleware((*Context).BotCheckMiddleware)
	}

	router.Middleware((*Context).SecurityHeadersMiddleware)
	router.Middleware((*Context).AuthMiddleware)

	router.Middleware((*Context).GlobalsMiddleware)
	router.Middleware((*Context).CurrencyMiddleware)

	// Images
	router.Get("/user-avatar/:user", (*Context).UserAvatar)
	router.Get("/item-category-image/:user", (*Context).ViewItemCategoryImage)

	router.Middleware((*Context).AbuseRateLimitMiddleware)
	router.Middleware((*Context).RateLimitMiddleware)

	// Index
	router.Get("/", (*Context).Index)
	router.Get("/captcha/:captcha_id", (*Context).ViewCaptchaImageV2)
	router.Get("botcheck", (*Context).BotCheckGET)
	router.Post("botcheck", (*Context).BotCheckPOST)

	loggedInRouter := router.Subrouter(Context{}, "/")
	loggedInRouter.Middleware((*Context).MessageStatsMiddleware)
	loggedInRouter.Middleware((*Context).TransactionStatsMiddleware)
	loggedInRouter.Middleware((*Context).DisputeStatsMiddleware)
	loggedInRouter.Middleware((*Context).WalletsMiddleware)

	// Static
	staticRouter := loggedInRouter.Subrouter(Context{}, "/help")
	staticRouter.Get("/", (*Context).Help)
	staticRouter.Get("/:filename", (*Context).HelpItem)

	// SERP
	loggedInRouter.Get("/marketplace/", (*Context).ListSerpItems)
	loggedInRouter.Get("/marketplace/:package_type", (*Context).ListSerpItems)
	loggedInRouter.Get("/vendors/", (*Context).ListSerpVendors)
	loggedInRouter.Get("/vendors/:package_type", (*Context).ListSerpVendors)

	// Auth
	authRouter := loggedInRouter.Subrouter(Context{}, "/auth")
	authRouter.Get("/login", (*Context).LoginGET)
	authRouter.Post("/login", (*Context).LoginPOST)
	authRouter.Get("/recover", (*Context).ViewRecoverGET)
	authRouter.Post("/recover", (*Context).ViewRecoverPOST)
	authRouter.Get("/register", (*Context).RegisterGET)
	authRouter.Post("/register", (*Context).RegisterPOST)
	authRouter.Get("/register/:invite_code", (*Context).RegisterGET)
	authRouter.Post("/register/:invite_code", (*Context).RegisterPOST)
	authRouter.Get("/mobile_auth", (*Context).ViewsAuthMobileAuthImage)
	authRouter.Post("/logout", (*Context).Logout)

	// Referral Admin
	referralAdminRouter := loggedInRouter.Subrouter(Context{}, "/referrals/admin")
	referralAdminRouter.Middleware((*Context).AdminMiddleware)
	referralAdminRouter.Get("/", (*Context).ViewAdminListReferralPayments)

	// Auth Admin
	authAdminRouter := authRouter.Subrouter(Context{}, "/admin")
	authAdminRouter.Middleware((*Context).AdminMiddleware)
	authAdminRouter.Get("/users", (*Context).AdminUsers)
	authAdminRouter.Post("/users/:user/login", (*Context).LoginAsUser)
	authAdminRouter.Post("/users/:user/ban", (*Context).BanUser)
	authAdminRouter.Post("/users/:user/staff", (*Context).GrantStaff)
	authAdminRouter.Post("/users/:user/seller", (*Context).GrantSeller)
	authAdminRouter.Get("/reviews", (*Context).AdminReviews)

	// Staff
	staffRouter := loggedInRouter.Subrouter(Context{}, "/staff")
	staffRouter.Middleware((*Context).StaffMiddleware)

	staffRouter.Get("", (*Context).ViewStaff)
	staffRouter.Get("/staff", (*Context).ViewStaffListStaff)
	staffRouter.Get("/users", (*Context).ViewStaffListSupportTickets)
	staffRouter.Get("/advertising", (*Context).ViewStaffAdvertisings)
	staffRouter.Get("/warnings", (*Context).ViewStaffListWarnings)
	staffRouter.Get("/messages", (*Context).ViewStaffListReportedMessages)
	staffRouter.Get("/messages/:message_uuid", (*Context).ViewStaffShowReportedMessage)
	staffRouter.Get("/items", (*Context).ViewStaffListItems)
	staffRouter.Get("/disputes", (*Context).ViewStaffListDisputes)
	staffRouter.Get("/stats", (*Context).ViewStaffStats)
	staffRouter.Get("/deposits", (*Context).ViewStaffListDeposits)

	// Staff - Stats
	staffRouter.Get("/stats/users.png", (*Context).ViewStatsNumberOfUsersGraph)
	staffRouter.Get("/stats/vendors.png", (*Context).ViewStatsNumberOfVendorsGraph)
	staffRouter.Get("/stats/trade-amount-btc.png", (*Context).ViewStatsBTCTradeAmountGraph)
	staffRouter.Get("/stats/trade-amount-eth.png", (*Context).ViewStatsETHTradeAmountGraph)
	staffRouter.Get("/stats/transactions.png", (*Context).ViewStatsNumberOfTransactionsGraph)

	// Staff - Stores
	staffRouter.Get("/vendors", (*Context).ViewStaffListVendors)
	staffRouter.Get("/vendors/:username", (*Context).ViewStaffVendorVerificationThreadGET)
	staffRouter.Post("/vendors/:username", (*Context).ViewStaffVendorVerificationThreadPOST)

	// Staff - User
	staffRouter.Get("/users/:username/payments", (*Context).ViewStaffUserPayments)
	staffRouter.Get("/users/:username/finance", (*Context).ViewStaffUserFinance)
	staffRouter.Get("/users/:username/tickets", (*Context).ViewStaffUserTickets)
	staffRouter.Get("/users/:username/tickets/:id", (*Context).ShowSupportTicketGET)
	staffRouter.Post("/users/:username/tickets/:id", (*Context).ShowSupportTicketPOST)
	staffRouter.Get("/users/:username/actions", (*Context).ViewStaffUserAdminActions)

	// Staff - User - Admin Actions
	staffRouter.Post("/users/:user/ban", (*Context).BanUser)
	staffRouter.Post("/users/:user/staff", (*Context).GrantStaff)
	staffRouter.Post("/users/:user/seller", (*Context).GrantSeller)

	// Staff - Store
	staffRouter.Get("/stores/:storename/payments", (*Context).ViewStaffStorePayments)
	staffRouter.Get("/stores/:storename/disputes", (*Context).ViewStaffStoreDisputes)
	staffRouter.Get("/stores/:storename/actions", (*Context).ViewStaffStoreAdminActions)

	// Staff - Store - Store Actions
	staffRouter.Post("/stores/:storename/suspend", (*Context).ViewStaffStoreToggleSuspend)
	staffRouter.Post("/stores/:storename/allow_to_sell", (*Context).ViewStaffStoreToggleAllowToSell)
	staffRouter.Post("/stores/:storename/trusted", (*Context).ViewStaffStoreToggleTrusted)

	staffRouter.Post("/stores/:storename/gold", (*Context).ViewStaffStoreToggleGoldAccount)
	staffRouter.Post("/stores/:storename/silver", (*Context).ViewStaffStoreToggleSilverAccount)
	staffRouter.Post("/stores/:storename/bronze", (*Context).ViewStaffStoreToggleBronzeAccount)
	staffRouter.Post("/stores/:storename/free", (*Context).ViewStaffStoreToggleFreeAccount)

	// Staff - CRUD
	staffRouter.Get("/item_categories", (*Context).ViewStaffCategories)
	staffRouter.Get("/item_categories/:id", (*Context).ViewStaffCategoriesEditGET)
	staffRouter.Post("/item_categories/:id", (*Context).ViewStaffCategoriesEditPOST)
	staffRouter.Post("/item_categories/:id/delete", (*Context).ViewStaffCategoriesDelete)

	// Warnings
	loggedInRouter.Get("/user/active_reservation", (*Context).ActiveReservation)

	// Store
	userRouter := loggedInRouter.Subrouter(Context{}, "/user/:username")
	userRouter.Middleware((*Context).UserMiddleware)
	userRouter.Get("/", (*Context).ViewAboutUser)

	// Store
	storeRouter := loggedInRouter.Subrouter(Context{}, "/store/:store")
	storeRouter.Middleware((*Context).PublicStoreMiddleware)
	storeRouter.Get("/", (*Context).ViewAboutStore)
	storeRouter.Get("/items", (*Context).ViewListStoreItems)
	storeRouter.Get("/reviews", (*Context).ViewStoreReviews)
	storeRouter.Get("/warnings", (*Context).ViewStoreWarningsGET)
	storeRouter.Post("/warnings", (*Context).ViewStoreWarningsPOST)
	storeRouter.Post("/warnings/:uuid", (*Context).ViewStoreWarningUpdateStatusPOST)

	// Store Item
	storeItemRouter := storeRouter.Subrouter(Context{}, "/item/:item")
	storeItemRouter.Middleware((*Context).PublicStoreItemMiddleware)
	storeItemRouter.Get("/", (*Context).ViewShowItem)
	storeItemRouter.Get("/package/:hash/book", (*Context).PreBookPackage)
	storeItemRouter.Post("/package/:hash/book", (*Context).BookPackage)

	// Board
	boardRouter := loggedInRouter.Subrouter(Context{}, "/board")
	boardRouter.Get("/message/:uuid/image", (*Context).MessageImage)
	boardRouter.Post("/:uuid/delete", (*Context).DeleteThread)
	boardRouter.Post("/:uuid/report", (*Context).ViewMessageReportPOST)

	// Messages
	messagesRouter := loggedInRouter.Subrouter(Context{}, "/messages")
	messagesRouter.Middleware((*Context).MessagesMiddleware)
	messagesRouter.Get("/", (*Context).ViewListPrivateMessages)
	messagesRouter.Get("/:username", (*Context).ViewShowPrivateMessageGET)
	messagesRouter.Post("/:username", (*Context).ViewShowPrivateMessagePOST)

	messagesAdminRouter := messagesRouter.Subrouter(Context{}, "/admin")
	messagesAdminRouter.Middleware((*Context).AdminMiddleware)
	messagesAdminRouter.Get("/:uuid", (*Context).AdminMessagesShow)

	// Package
	packagesRouter := loggedInRouter.Subrouter(Context{}, "/packages")
	packagesRouter.Get("/", (*Context).ListPackages)

	singlePackageRouter := packagesRouter.Subrouter(Context{}, "/:package")
	singlePackageRouter.Middleware((*Context).PackageMiddleware)
	singlePackageRouter.Get("/image", (*Context).PackageImage)

	// Profile
	loggedInRouter.Get("/settings/user", (*Context).ViewUserSettingsGET)
	loggedInRouter.Post("/settings/user", (*Context).ViewUserSettingsPOST)
	loggedInRouter.Get("/settings/store", (*Context).ViewStoreSettingsGET)
	loggedInRouter.Post("/settings/store", (*Context).ViewStoreSettingsPOST)
	loggedInRouter.Get("/referrals", (*Context).Referrals)
	loggedInRouter.Get("/verification/encryption", (*Context).ViewVerificationEncryptionGET)
	loggedInRouter.Get("/verification/agreement", (*Context).ViewVerificationAgreementGET)
	loggedInRouter.Post("/verification/agreement", (*Context).ViewVerificationAgreementPOST)
	loggedInRouter.Get("/verification/plan", (*Context).ViewVerificationPlanGET)
	loggedInRouter.Get("/verification/plan/:account", (*Context).ViewVerificationPlanPurchaseGET)
	loggedInRouter.Post("/verification/plan/:account", (*Context).ViewVerificationPlanPurchasePOST)
	loggedInRouter.Get("/settings/user/pgp/step1", (*Context).SetupPGPViaDecryptionStep1GET)
	loggedInRouter.Post("/settings/user/pgp/step1", (*Context).SetupPGPViaDecryptionStep1POST)
	loggedInRouter.Post("/settings/user/pgp/step2", (*Context).SetupPGPViaDecryptionStep2POST)
	loggedInRouter.Get("/settings/currency/:currency", (*Context).SetCurrency)
	loggedInRouter.Get("/settings/language/:language", (*Context).SetLanguage)

	// Profile
	loggedInRouter.Post("/shipping", (*Context).SaveShippingOption)
	loggedInRouter.Post("/shipping/delete", (*Context).DeleteShippingOption)

	// Support
	supportRouter := loggedInRouter.Subrouter(Context{}, "/support")
	supportRouter.Get("/", (*Context).ViewListSupportTickets)
	supportRouter.Get("/new", (*Context).ShowSupportTicketGET)
	supportRouter.Post("/new", (*Context).ShowSupportTicketPOST)
	supportRouter.Get("/:id", (*Context).ShowSupportTicketGET)
	supportRouter.Post("/:id", (*Context).ShowSupportTicketPOST)
	supportRouter.Post("/:id/status", (*Context).ViewUpdateTicketStatusPOST)

	// Store CMS
	sellerRouter := loggedInRouter.Subrouter(Context{}, "/store-admin/:store")
	sellerRouter.Middleware((*Context).PrivateStoreMiddleware)

	// Advertisings
	sellerRouter.Get("/advertisings", (*Context).ViewListAdvertisings)
	sellerRouter.Post("/advertisings", (*Context).AddAdvertisingsPOST)

	// Deposit
	depositRouter := sellerRouter.Subrouter(Context{}, "/deposits")
	depositRouter.Get("/", (*Context).ViewShowDeposit)
	depositRouter.Get("/add", (*Context).ViewShowDepositAddGET)
	depositRouter.Post("/add", (*Context).ViewShowDepositAddPOST)
	depositRouter.Get("/:deposit_uuid/withdraw", (*Context).ViewWithdrawDepositGET)
	depositRouter.Post("/:deposit_uuid/withdraw", (*Context).ViewWithdrawDepositPOST)

	// Items CMS
	itemRouter := sellerRouter.Subrouter(Context{}, "/item/:item")
	itemRouter.Middleware((*Context).PrivateStoreItemMiddleware)
	itemRouter.Get("/edit", (*Context).EditItem)
	itemRouter.Post("/edit", (*Context).SaveItem)
	itemRouter.Post("/delete", (*Context).DeleteItem)

	// Package CMS
	packageRouter := itemRouter.Subrouter(Context{}, "/package/:package")
	packageRouter.Middleware((*Context).PrivateStoreItemPackageMiddleware)
	packageRouter.Get("/edit", (*Context).EditPackageStep1)
	packageRouter.Post("/edit", (*Context).SavePackage)
	packageRouter.Post("/duplicate", (*Context).DuplicatePackage)
	packageRouter.Post("/delete", (*Context).DeletePackage)

	// Crowdloans
	crowdloansRouter := loggedInRouter.Subrouter(Context{}, "/crowdloans")
	crowdloansRouter.Middleware((*Context).PolkadotWalletMiddleware)
	crowdloansRouter.Get("/", (*Context).ListCrowdloans)
	crowdloansRouter.Get("/mint", (*Context).MintCrowdloan)
	crowdloansRouter.Post("/mint", (*Context).MintCrowdloanPOST)
	crowdloansRouter.Get("/loans", (*Context).ListUserLoans)
	crowdloansRouter.Get("/lends", (*Context).ListUserLends)

	singleCrowdloanRouter := crowdloansRouter.Subrouter(Context{}, "/:crowdloan")
	singleCrowdloanRouter.Middleware((*Context).CrowdloanMiddleware)
	singleCrowdloanRouter.Get("/", (*Context).ShowCrowdloan)
	singleCrowdloanRouter.Post("/fund", (*Context).FundCrowdloanPOST)
	singleCrowdloanRouter.Post("/withdraw", (*Context).WithdrawCrowdloanPOST)
	singleCrowdloanRouter.Post("/payback", (*Context).PaybackCrowdloanPOST)
	singleCrowdloanRouter.Post("/payout", (*Context).PayoutCrowdloanPOST)

	// Transactions
	transactionRouter := loggedInRouter.Subrouter(Context{}, "/payments")
	transactionRouter.Get("/", (*Context).ListCurrentTransactionStatuses)

	singleTransactionRouter := transactionRouter.Subrouter(Context{}, "/:transaction")
	singleTransactionRouter.Middleware((*Context).TransactionMiddleware)
	singleTransactionRouter.Get("/", (*Context).ShowTransaction)
	singleTransactionRouter.Get("/image", (*Context).TransactionImage)
	singleTransactionRouter.Post("/", (*Context).ShowTransactionPOST)
	singleTransactionRouter.Post("/shipping", (*Context).SetTransactionShippingStatus)
	singleTransactionRouter.Post("/release", (*Context).ReleaseTransaction)
	singleTransactionRouter.Post("/cancel", (*Context).CancelTransaction)
	singleTransactionRouter.Post("/complete", (*Context).CompleteTransactionPOST)

	// Disputes
	disputeRouter := loggedInRouter.Subrouter(Context{}, "/dispute")
	disputeRouter.Get("/", (*Context).ViewDisputeList)
	disputeRouter.Post("/start/:transaction_uuid", (*Context).ViewDisputeStart)

	concreteDisputeRouter := disputeRouter.Subrouter(Context{}, "/:dispute_uuid")
	concreteDisputeRouter.Middleware((*Context).DisputeMiddleware)
	concreteDisputeRouter.Get("/", (*Context).ViewDispute)
	concreteDisputeRouter.Post("/status", (*Context).ViewDisputeSetStatus)
	concreteDisputeRouter.Post("/claim", (*Context).ViewDisputeClaimAdd)
	// concreteDisputeRouter.Post("/partial_refund", (*Context).ViewDisputePartialRefund)

	disputeClaimRouter := concreteDisputeRouter.Subrouter(Context{}, "/:dispute_claim_id")
	disputeClaimRouter.Middleware((*Context).DisputeClaimMiddleware)
	disputeClaimRouter.Get("/", (*Context).ViewDisputeClaimGET)
	disputeClaimRouter.Post("/", (*Context).ViewDisputeClaimPOST)

	disputeAdminRouter := disputeRouter.Subrouter(Context{}, "/admin")
	disputeAdminRouter.Middleware((*Context).AdminMiddleware)
	disputeAdminRouter.Get("/list", (*Context).AdminDisputeList)

	// Wallet
	walletRouter := loggedInRouter.Subrouter(Context{}, "/wallet")

	// Bitcoin Wallet
	bitcoinWalletRouter := walletRouter.Subrouter(Context{}, "/bitcoin")
	bitcoinWalletRouter.Middleware((*Context).BitcoinWalletMiddleware)
	bitcoinWalletRouter.Get("/receive", (*Context).BitcoinWalletRecieve)
	bitcoinWalletRouter.Get("/send", (*Context).BitcoinWalletSendStep1GET)
	bitcoinWalletRouter.Post("/send", (*Context).BitcoinWalletSendPOST)
	bitcoinWalletRouter.Get("/:address/image", (*Context).BitcoinWalletImage)
	bitcoinWalletRouter.Get("/actions", (*Context).BitcoinWalletActions)

	// Ethereum Wallet
	ethereumWalletRouter := walletRouter.Subrouter(Context{}, "/ethereum")
	ethereumWalletRouter.Middleware((*Context).EthereumWalletMiddleware)
	ethereumWalletRouter.Get("/receive", (*Context).EthereumWalletRecieve)
	ethereumWalletRouter.Get("/send", (*Context).EthereumWalletSendGET)
	ethereumWalletRouter.Post("/send", (*Context).EthereumWalletSendPOST)
	ethereumWalletRouter.Get("/:address/image", (*Context).EthereumWalletImage)
	ethereumWalletRouter.Get("/actions", (*Context).EthereumWalletActions)

	// Monero Wallet
	moneroWalletRouter := walletRouter.Subrouter(Context{}, "/monero")
	moneroWalletRouter.Middleware((*Context).MoneroWalletMiddleware)
	moneroWalletRouter.Get("/receive", (*Context).MoneroWalletRecieve)
	moneroWalletRouter.Get("/send", (*Context).MoneroWalletSendGET)
	moneroWalletRouter.Post("/send", (*Context).MoneroWalletSendPOST)
	moneroWalletRouter.Get("/:address/image", (*Context).MoneroWalletImage)
	moneroWalletRouter.Get("/actions", (*Context).MoneroWalletActions)

	// Polkadot Wallet
	polkadotWalletRouter := walletRouter.Subrouter(Context{}, "/polkadot")
	polkadotWalletRouter.Middleware((*Context).PolkadotWalletMiddleware)
	polkadotWalletRouter.Get("/receive", (*Context).PolkadotWalletRecieve)
	polkadotWalletRouter.Get("/send", (*Context).PolkadotWalletSendGET)
	polkadotWalletRouter.Post("/send", (*Context).PolkadotWalletSendPOST)
	polkadotWalletRouter.Get("/:address/image", (*Context).PolkadotWalletImage)
	polkadotWalletRouter.Get("/actions", (*Context).PolkadotWalletActions)
	polkadotWalletRouter.Get("/mnemonic", (*Context).PolkadotWalletMnemonic)

	// Transactions Admin
	transactionAdminRouter := transactionRouter.Subrouter(Context{}, "/admin")
	transactionAdminRouter.Middleware((*Context).AdminMiddleware)
	transactionAdminRouter.Get("/list", (*Context).AdminListTransactions)
	transactionAdminRouter.Post("/:transaction/cancel", (*Context).AdminTransactionCancel)
	transactionAdminRouter.Post("/:transaction/fail", (*Context).AdminTransactionFail)
	transactionAdminRouter.Post("/:transaction/pending", (*Context).AdminTransactionMakePending)
	transactionAdminRouter.Post("/:transaction/complete", (*Context).AdminTransactionComplete)
	transactionAdminRouter.Post("/:transaction/release", (*Context).AdminTransactionRelease)
	transactionAdminRouter.Post("/:transaction/freeze", (*Context).AdminTransactionFreeze)
	transactionAdminRouter.Post("/:transaction/update", (*Context).AdminTransactionUpdateStatus)

	return router
}
