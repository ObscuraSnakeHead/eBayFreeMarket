package webapp

import (
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"time"

	"github.com/gocraft/web"

	"ucoi4tzlito52agmquc6oopn3zpmr6djz5vvfabtgrhyc6hufpzjtnad.onion/Tochka/tochka-free-market/modules/util"
)

func (c *Context) ListCrowdloans(w web.ResponseWriter, r *web.Request) {
	if len(r.URL.Query()["status"]) > 0 {
		c.SelectedStatus = r.URL.Query()["status"][0]
	}

	viewCrowdloans := FindAllCrowdloans().ViewCrowdloans(c.ViewUser.Language)
	c.ViewCrowdloans = viewCrowdloans

	// filter crowdloans by STATUS
	if c.SelectedStatus != "" {
		filteredViewCrowdloans := ViewCrowdloans{}
		for _, viewCrowdloan := range c.ViewCrowdloans {
			if viewCrowdloan.CurrentStatus == c.SelectedStatus {
				filteredViewCrowdloans = append(filteredViewCrowdloans, viewCrowdloan)
			}
		}
		c.ViewCrowdloans = filteredViewCrowdloans
	}

	util.RenderTemplate(w, "crowdloan/list", c)
}

func (c *Context) MintCrowdloan(w web.ResponseWriter, r *web.Request) {
	util.RenderTemplate(w, "crowdloan/mint", c)
}

func (c *Context) ListUserLoans(w web.ResponseWriter, r *web.Request) {
	loans, err := c.UserPolkadotWallet.MintedCrowdloans()
	if err != nil {
		http.NotFound(w, r.Request)
		return
	}

	dbCrowdloans := Crowdloans{}
	for _, crowdloan := range loans {
		dbCrowdloan := ConvertPolkadotCrowdloanFromBlockchainToModel(crowdloan, *c.ViewUser.User)
		dbCrowdloans = append(dbCrowdloans, dbCrowdloan)
	}

	c.ViewCrowdloans = dbCrowdloans.ViewCrowdloans(c.ViewUser.Language)

	sort.Slice(c.ViewCrowdloans, func(i, j int) bool {
		return c.ViewCrowdloans[i].PolkadotCrowdloanId > c.ViewCrowdloans[j].PolkadotCrowdloanId
	})

	util.RenderTemplate(w, "crowdloan/list_user_loans", c)
}

func (c *Context) ListUserLends(w web.ResponseWriter, r *web.Request) {
	loans, err := c.UserPolkadotWallet.FundedCrowdloans()
	if err != nil {
		http.NotFound(w, r.Request)
		return
	}

	for _, loan := range loans {
		viewCrowdloan := ConvertPolkadotCrowdloanFromBlockchainToModel(loan.Loan, *c.ViewUser.User).
			ViewCrowdloan(c.ViewUser.Language)
		c.ViewCrowdloanLends = append(c.ViewCrowdloanLends, ViewCrowdloanLend{
			ViewCrowdloan: &viewCrowdloan,
			LendValue:     float64(loan.LendValue),
		})
	}

	sort.Slice(c.ViewCrowdloans, func(i, j int) bool {
		return c.ViewCrowdloans[i].PolkadotCrowdloanId > c.ViewCrowdloans[j].PolkadotCrowdloanId
	})

	util.RenderTemplate(w, "crowdloan/list_user_lends", c)
}

func (c *Context) MintCrowdloanPOST(w web.ResponseWriter, r *web.Request) {
	goalValue, err := strconv.ParseFloat(r.FormValue("amount"), 64)
	if err != nil {
		c.Error = err.Error()
		c.MintCrowdloan(w, r)
		return
	}

	durationWeeks, err := strconv.ParseInt(r.FormValue("duration"), 10, 64)
	if err != nil {
		c.Error = err.Error()
		c.MintCrowdloan(w, r)
		return
	}

	weeklyInterest, err := strconv.ParseInt(r.FormValue("weekly_interest"), 10, 64)
	if err != nil {
		c.Error = err.Error()
		c.MintCrowdloan(w, r)
		return
	}

	now := time.Now()
	goalTimestamp := now.Add(time.Duration(durationWeeks) * 7 * 24 * time.Hour)

	crowdloan, err := c.UserPolkadotWallet.MintCrowdloan(
		goalValue,
		goalTimestamp,
		uint16(weeklyInterest),
	)
	if err != nil {
		c.Error = err.Error()
		c.MintCrowdloan(w, r)
		return
	}

	url := fmt.Sprintf("/crowdloans/polkadot-%d", crowdloan.PolkadotCrowdloanId)
	http.Redirect(w, r.Request, url, http.StatusFound)
}

func (c *Context) ShowCrowdloan(w web.ResponseWriter, r *web.Request) {
	lends, err := c.UserPolkadotWallet.FundedCrowdloans()
	if err != nil {
		http.NotFound(w, r.Request)
		return
	}

	for _, lend := range lends {
		if lend.CrowdloanId == c.ViewCrowdloan.PolkadotCrowdloanId {
			viewCrowdloan := ConvertPolkadotCrowdloanFromBlockchainToModel(lend.Loan, *c.ViewUser.User).
				ViewCrowdloan(c.ViewUser.Language)
			c.ViewCrowdloanLend = &ViewCrowdloanLend{
				ViewCrowdloan: &viewCrowdloan,
				LendValue:     float64(lend.LendValue),
			}

			break
		}
	}

	util.RenderTemplate(w, "crowdloan/show", c)
}

func (c *Context) FundCrowdloanPOST(w web.ResponseWriter, r *web.Request) {
	amount, err := strconv.ParseFloat(r.FormValue("amount"), 64)
	if err != nil {
		c.ShowCrowdloan(w, r)
		return
	}

	_, err = c.UserPolkadotWallet.FundCrowdloan(c.ViewCrowdloan.PolkadotCrowdloanId, amount)
	if err != nil {
		c.Error = err.Error()
		c.MintCrowdloan(w, r)
		return
	}

	url := fmt.Sprintf("/crowdloans/polkadot-%d", c.ViewCrowdloan.PolkadotCrowdloanId)
	http.Redirect(w, r.Request, url, http.StatusFound)
}

func (c *Context) WithdrawCrowdloanPOST(w web.ResponseWriter, r *web.Request) {

	_, err := c.UserPolkadotWallet.WithdrawCrowdloan(c.ViewCrowdloan.PolkadotCrowdloanId)
	if err != nil {
		c.Error = err.Error()
		c.MintCrowdloan(w, r)
		return
	}

	url := fmt.Sprintf("/crowdloans/polkadot-%d", c.ViewCrowdloan.PolkadotCrowdloanId)
	http.Redirect(w, r.Request, url, http.StatusFound)
}

func (c *Context) PaybackCrowdloanPOST(w web.ResponseWriter, r *web.Request) {

	_, err := c.UserPolkadotWallet.PaybackCrowdloan(c.ViewCrowdloan.PolkadotCrowdloanId)
	if err != nil {
		c.Error = err.Error()
		c.MintCrowdloan(w, r)
		return
	}

	url := fmt.Sprintf("/crowdloans/polkadot-%d", c.ViewCrowdloan.PolkadotCrowdloanId)
	http.Redirect(w, r.Request, url, http.StatusFound)
}

func (c *Context) PayoutCrowdloanPOST(w web.ResponseWriter, r *web.Request) {

	_, err := c.UserPolkadotWallet.PayoutCrowdloan(c.ViewCrowdloan.PolkadotCrowdloanId)
	if err != nil {
		c.Error = err.Error()
		c.MintCrowdloan(w, r)
		return
	}

	url := fmt.Sprintf("/crowdloans/polkadot-%d", c.ViewCrowdloan.PolkadotCrowdloanId)
	http.Redirect(w, r.Request, url, http.StatusFound)
}
