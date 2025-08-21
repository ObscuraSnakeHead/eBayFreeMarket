package webapp

import (
	"github.com/gocraft/web"

	"ucoi4tzlito52agmquc6oopn3zpmr6djz5vvfabtgrhyc6hufpzjtnad.onion/Tochka/tochka-free-market/modules/util"
)

func (c *Context) ViewAboutUser(w web.ResponseWriter, r *web.Request) {
	c.SelectedSection = "info"
	if len(r.URL.Query()["section"]) > 0 {
		c.SelectedSection = r.URL.Query()["section"][0]
	}
	util.RenderTemplateOrAPIResponse(w, r, "user/about", c, c.IsAPIRequest)
}
