package webapp

import (
	"net/http"

	"github.com/gocraft/web"
)

func (c *Context) Index(w web.ResponseWriter, r *web.Request) {
	redirectUrl := "/marketplace"
	http.Redirect(w, r.Request, redirectUrl, http.StatusFound)
}
