package fix

import (
	"fmt"
	"ipgw/base/ctx"
	"net/http"
)

func fix() {
	fmt.Println(fixing)
	x := ctx.GetCtx()
	x.User.Cookie = &http.Cookie{}
	x.SaveAll()
	fmt.Println(successFix)
}
