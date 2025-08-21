package webapp

import (
	"net/http"
	"strconv"

	btcqr "github.com/GeertJohan/go.btcqr"
	"github.com/dchest/captcha"
	"github.com/gocraft/web"
	"github.com/mojocn/base64Captcha"
	"ucoi4tzlito52agmquc6oopn3zpmr6djz5vvfabtgrhyc6hufpzjtnad.onion/Tochka/tochka-free-market/modules/util"
)

func (c *Context) PolkadotWalletRecieve(w web.ResponseWriter, r *web.Request) {
	util.RenderTemplate(w, "wallet/polkadot/recieve", c)
}

func (c *Context) PolkadotWalletSendGET(w web.ResponseWriter, r *web.Request) {
	c.CaptchaId = captcha.New()
	util.RenderTemplate(w, "wallet/polkadot/send", c)
}

func (c *Context) PolkadotWalletSendPOST(w web.ResponseWriter, r *web.Request) {
	isCaptchaValid := base64Captcha.VerifyCaptcha(r.FormValue("captcha_id"), r.FormValue("captcha"))
	if !isCaptchaValid {
		c.Error = "Invalid captcha"
		c.PolkadotWalletSendGET(w, r)
		return
	}

	var (
		address   = r.FormValue("address")
		amountStr = r.FormValue("amount")
	)

	if !polkadotRegexp.MatchString(address) {
		c.Error = "Wrong DOT address"
		c.PolkadotWalletSendGET(w, r)
		return
	}

	amount, err := strconv.ParseFloat(amountStr, 64)
	if err != nil {
		c.Error = "Wrong amount"
		c.PolkadotWalletSendGET(w, r)
		return
	}

	results, err := c.UserPolkadotWallet.Send(address, amount)
	if err != nil {
		c.Error = err.Error()
		c.PolkadotWalletSendGET(w, r)
		return
	}

	dotPaymentResult, err := results.DOTPaymentResult()
	if err == nil {
		c.DOTPaymentResult = &dotPaymentResult
	}
	util.RenderTemplate(w, "wallet/polkadot/send_receipt", c)
}

func (c *Context) PolkadotWalletImage(w web.ResponseWriter, r *web.Request) {
	req := &btcqr.Request{
		Address: r.PathParams["address"],
	}
	code, err := req.GenerateQR()
	if err != nil {
		http.NotFound(w, r.Request)
		return
	}
	png := code.PNG()
	w.Header().Set("Content-type", "image/png")
	w.Write(png)
}

func (c *Context) PolkadotWalletActions(w web.ResponseWriter, r *web.Request) {
	c.UserPolkadotWalletActions = FindUserPolkadotWalletActionsForUser(c.ViewUser.Uuid)
	util.RenderTemplateOrAPIResponse(w, r, "wallet/polkadot/actions", c, c.IsAPIRequest)
}

func (c *Context) PolkadotWalletMnemonic(w web.ResponseWriter, r *web.Request) {
	util.RenderTemplateOrAPIResponse(w, r, "wallet/polkadot/mnemonic", c, c.IsAPIRequest)
}
