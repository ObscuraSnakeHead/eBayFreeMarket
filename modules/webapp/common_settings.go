package webapp

import (
	"regexp"

	"github.com/microcosm-cc/bluemonday"
	"ucoi4tzlito52agmquc6oopn3zpmr6djz5vvfabtgrhyc6hufpzjtnad.onion/Tochka/tochka-free-market/modules/settings"
)

/*
	Globals
*/

var (
	bitcoinRegexp      = regexp.MustCompile("^(bc(0([ac-hj-np-z02-9]{39}|[ac-hj-np-z02-9]{59})|1[ac-hj-np-z02-9]{8,87})|[13][a-km-zA-HJ-NP-Z1-9]{25,35})$")
	emailRegexp        = regexp.MustCompile("^[a-zA-Z0-9_.+-]+@[a-zA-Z0-9-]+\\.[a-zA-Z0-9-.]+$")
	ethereumRegexp     = regexp.MustCompile("^0x[a-fA-F0-9]{40}$")
	moneroRegexp       = regexp.MustCompile("^[48][0-9AB][1-9A-HJ-NP-Za-km-z]{93}$")
	openBitcoinRegexp  = regexp.MustCompile("[13][a-km-zA-HJ-NP-Z1-9]{25,34}")
	openEthereumRegexp = regexp.MustCompile("0x[a-fA-F0-9]{40}")
	polkadotRegexp     = regexp.MustCompile("^[1-9A-HJ-NP-Za-km-z]{47,48}$")
	usernameRegexp     = regexp.MustCompile("^[a-z0-9_-]{3,16}$")

	userHtmlPolicy         = bluemonday.NewPolicy()
	htmlPolicy             = bluemonday.UGCPolicy()
	messageboardHtmlPolicy = bluemonday.NewPolicy()

	MARKETPLACE_SETTINGS = settings.GetSettings()
)

func init() {
	userHtmlPolicy.AllowElements("h1", "h2", "h3", "h4", "h5", "p", "strong", "i")
}
