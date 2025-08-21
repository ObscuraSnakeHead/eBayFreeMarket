package webapp

import (
	"net/http"

	"github.com/gocraft/web"
)

func (c *Context) CrowdloanMiddleware(w web.ResponseWriter, r *web.Request, next web.NextMiddlewareFunc) {
	crowdloan, _ := FindCrowdloanByUuid(r.PathParams["crowdloan"])
	if crowdloan == nil {
		http.NotFound(w, r.Request)
		return
	}

	err := crowdloan.UpdateStatus()
	if err != nil {
		http.NotFound(w, r.Request)
		return
	}

	viewCrowdloan := crowdloan.ViewCrowdloan(c.ViewUser.Language)
	c.ViewCrowdloan = &viewCrowdloan

	next(w, r)
}
