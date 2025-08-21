package webapp

import (
	"github.com/gocraft/web"
	"ucoi4tzlito52agmquc6oopn3zpmr6djz5vvfabtgrhyc6hufpzjtnad.onion/Tochka/tochka-free-market/modules/apis"
)

func (c *Context) BitcoinWalletMiddleware(w web.ResponseWriter, r *web.Request, next web.NextMiddlewareFunc) {
	c.UserBitcoinWallets = c.ViewUser.FindUserBitcoinWallets()
	for _, w := range c.UserBitcoinWallets {
		w.UpdateBalance(false)
	}
	if len(c.UserBitcoinWallets) > 0 {
		c.UserBitcoinWallet = &c.UserBitcoinWallets[0]
	} else {
		c.UserBitcoinWallet, _ = CreateBitcoinWallet(*c.ViewUser.User)
		c.UserBitcoinWallets = append(c.UserBitcoinWallets, *c.UserBitcoinWallet)
	}
	balance := c.UserBitcoinWallets.Balance()
	c.UserBitcoinBalance = &balance
	next(w, r)
}

func (c *Context) EthereumWalletMiddleware(w web.ResponseWriter, r *web.Request, next web.NextMiddlewareFunc) {
	c.UserEthereumWallets = c.ViewUser.FindUserEthereumWallets()
	for _, w := range c.UserEthereumWallets {
		w.UpdateBalance(false)
	}
	if len(c.UserEthereumWallets) > 0 {
		c.UserEthereumWallet = &c.UserEthereumWallets[0]
	} else {
		c.UserEthereumWallet, _ = CreateEthereumWallet(*c.ViewUser.User)
		c.UserEthereumWallets = append(c.UserEthereumWallets, *c.UserEthereumWallet)
	}
	c.UserEthereumBalance = &apis.ETHWalletBalance{
		Balance: c.UserEthereumWallets.Balance().Balance,
	}
	next(w, r)
}

func (c *Context) MoneroWalletMiddleware(w web.ResponseWriter, r *web.Request, next web.NextMiddlewareFunc) {
	c.UserMoneroWallets = c.ViewUser.FindUserMoneroWallets()
	for _, w := range c.UserMoneroWallets {
		w.UpdateBalance(false)
	}
	if len(c.UserMoneroWallets) > 0 {
		c.UserMoneroWallet = &c.UserMoneroWallets[0]
	} else {
		c.UserMoneroWallet, _ = CreateMoneroWallet(*c.ViewUser.User)
		c.UserMoneroWallets = append(c.UserMoneroWallets, *c.UserMoneroWallet)
	}
	c.UserMoneroBalance = &apis.XMRWalletBalance{
		Balance: c.UserMoneroWallets.Balance().Balance,
	}
	next(w, r)
}

func (c *Context) PolkadotWalletMiddleware(w web.ResponseWriter, r *web.Request, next web.NextMiddlewareFunc) {
	c.UserPolkadotWallets = c.ViewUser.FindUserPolkadotWallets()
	for _, w := range c.UserPolkadotWallets {
		w.UpdateBalance(false)
	}
	if len(c.UserPolkadotWallets) > 0 {
		c.UserPolkadotWallet = &c.UserPolkadotWallets[0]
	} else {
		c.UserPolkadotWallet, _ = CreatePolkadotWallet(*c.ViewUser.User)
		c.UserPolkadotWallets = append(c.UserPolkadotWallets, *c.UserPolkadotWallet)
	}
	dotBalance := c.UserPolkadotWallets.Balance()
	c.UserPolkadotBalance = &apis.DOTWalletBalance{
		FreeBalance:     dotBalance.FreeBalance,
		ReservedBalance: dotBalance.ReservedBalance,
		FrozenBalance:   dotBalance.FrozenBalance,
	}

	next(w, r)
}

func (c *Context) WalletsMiddleware(w web.ResponseWriter, r *web.Request, next web.NextMiddlewareFunc) {
	if c.ViewUser != nil {
		// ETH
		c.UserEthereumWallets = c.ViewUser.FindUserEthereumWallets()
		c.UserEthereumBalance = &apis.ETHWalletBalance{
			Balance: c.UserEthereumWallets.Balance().Balance,
		}

		// BTC
		c.UserBitcoinWallets = c.ViewUser.FindUserBitcoinWallets()
		btcBalance := c.UserBitcoinWallets.Balance()
		c.UserBitcoinBalance = &btcBalance

		// XMR
		c.UserMoneroWallets = c.ViewUser.FindUserMoneroWallets()
		xmrBalance := c.UserMoneroWallets.Balance()
		c.UserMoneroBalance = &xmrBalance

		// DOT
		c.UserPolkadotWallets = c.ViewUser.FindUserPolkadotWallets()
		dotBalance := c.UserPolkadotWallets.Balance()
		c.UserPolkadotBalance = &dotBalance
	}

	next(w, r)
}
