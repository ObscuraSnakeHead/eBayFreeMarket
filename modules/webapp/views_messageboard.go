package webapp

import (
	"net/http"

	"github.com/gocraft/web"

	"ucoi4tzlito52agmquc6oopn3zpmr6djz5vvfabtgrhyc6hufpzjtnad.onion/Tochka/tochka-free-market/modules/util"
)

func (c *Context) MessageImage(w web.ResponseWriter, r *web.Request) {
	size := "normal"
	if len(r.URL.Query()["size"]) > 0 {
		size = r.URL.Query()["size"][0]
	}
	util.ServeImage(r.PathParams["uuid"], size, w, r)
}

func (c *Context) DeleteThread(w web.ResponseWriter, r *web.Request) {
	thread, err := FindThreadByUuid(r.PathParams["uuid"])
	if err != nil {
		http.NotFound(w, r.Request)
		return
	}
	if thread.SenderUserUuid == c.ViewUser.Uuid || c.ViewUser.IsAdmin || c.ViewUser.IsStaff {
		thread.Remove()
	}
	http.Redirect(w, r.Request, "/board/", http.StatusFound)
}

func (c *Context) ViewMessageReportPOST(w web.ResponseWriter, r *web.Request) {
	message, err := FindMessageByUuid(r.PathParams["uuid"])
	if err != nil {
		http.NotFound(w, r.Request)
		return
	}
	if message.RecieverUserUuid != c.ViewUser.Uuid {
		http.NotFound(w, r.Request)
		return
	}
	message.IsReported = true
	message.Save()
	redirectUrl := "/messages/" + message.SenderUser.Username
	http.Redirect(w, r.Request, redirectUrl, http.StatusFound)
}
