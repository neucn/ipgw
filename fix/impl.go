package fix

import (
	. "ipgw/base"
	"ipgw/ctx"
)

func fix() {
	InfoL(fixing)
	c := ctx.NewCtx()
	c.SaveAll()
	InfoL(successFix)
}
