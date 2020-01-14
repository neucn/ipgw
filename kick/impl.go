package kick

import (
	. "ipgw/base"
	"ipgw/core/gw"
	"ipgw/ctx"
)

func kick(c *ctx.Ctx, sid string) {
	InfoF(infoBeginKick, sid)

	// todo 这里也因为封装而产生了问题，遇到网络错误会直接结束程序而不会重试
	ok := gw.Kick(c, sid)

	if !ok {
		ErrorF(failKick, sid)
		return
	}

	InfoF(successKick, sid)
}
