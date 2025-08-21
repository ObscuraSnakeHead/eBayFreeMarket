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

func (c *Context) MoneroWalletRecieve(w web.ResponseWriter, r *web.Request) {
	util.RenderTemplate(w, "wallet/monero/recieve", c)
}

func (c *Context) MoneroWalletSendGET(w web.ResponseWriter, r *web.Request) {
	c.CaptchaId = captcha.New()
	util.RenderTemplate(w, "wallet/monero/send", c)
}

func (c *Context) MoneroWalletSendPOST(w web.ResponseWriter, r *web.Request) {
	isCaptchaValid := base64Captcha.VerifyCaptcha(r.FormValue("captcha_id"), r.FormValue("captcha"))
	if !isCaptchaValid {
		c.Error = "Invalid captcha"
		c.MoneroWalletSendGET(w, r)
		return
	}

	var (
		address   = r.FormValue("address")
		amountStr = r.FormValue("amount")
	)

	if !moneroRegexp.MatchString(address) {
		c.Error = "Wrong XMR address"
		c.MoneroWalletSendGET(w, r)
		return
	}

	amount, err := strconv.ParseFloat(amountStr, 64)
	if err != nil {
		c.Error = "Wrong amount"
		c.MoneroWalletSendGET(w, r)
		return
	}

	results, err := c.UserMoneroWallet.Send(address, amount)
	if err != nil {
		c.Error = err.Error()
		c.MoneroWalletSendGET(w, r)
		return
	}

	xmrPaymentResult, err := results.XMRPaymentResult()
	if err == nil {
		c.XMRPaymentResult = xmrPaymentResult
	}
	util.RenderTemplate(w, "wallet/monero/send_receipt", c)
}

func (c *Context) MoneroWalletImage(w web.ResponseWriter, r *web.Request) {
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

func (c *Context) MoneroWalletActions(w web.ResponseWriter, r *web.Request) {
	c.UserMoneroWalletActions = FindUserMoneroWalletActionsForUser(c.ViewUser.Uuid)
	util.RenderTemplateOrAPIResponse(w, r, "wallet/monero/actions", c, c.IsAPIRequest)
}
