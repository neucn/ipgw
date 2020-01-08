package fix

import (
	"fmt"
	"ipgw/ctx"
)

func fix() {
	fmt.Println(fixing)
	c := ctx.NewCtx()
	c.SaveAll()
	fmt.Println(successFix)
}
