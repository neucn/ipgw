package logout

import (
	. "ipgw/base"
	"ipgw/core/cas"
	. "ipgw/core/global"
	"ipgw/core/gw"
	"ipgw/ctx"
	"net/http"
	"net/url"
)

// 使用文件中保存的SID登出
// 这个SID是每次成功登陆都会保存，登出时并不含删除
// 由于按照层级关系，无参时默认 SID > Cookie > UP，所以失败时不能直接Fatal，而应该return false
func logoutWithSID(c *ctx.Ctx) (ok bool) {
	// 判断是否保存有SID
	if c.Net.SID == "" {
		return false
	}

	if ctx.FullView {
		InfoF(usingSID, c.Net.SID)
	}

	// 强制SID下线
	success := gw.Kick(c, c.Net.SID)

	if !success {
		if ctx.FullView {
			ErrorF(failLogoutBySID, c.Net.SID)
		}
		return false
	}

	InfoL(successLogoutBySID)
	return true
}

// 使用账号登出，其实就是先登录再登出
// 因为使用账号登陆是优先级最低的，因此失败时直接结束程序，不需要返回bool
func logoutWithUP(c *ctx.Ctx) {
	if ctx.FullView {
		InfoF(usingUP, c.User.Username)
	}

	reqUrl := "https://pass.neu.edu.cn/tpass/login?service=https%3A%2F%2Fipgw.neu.edu.cn%2Fsrun_cas.php%3Fac_id%3D1"

	// 获取必要参数
	lt, postUrl := cas.GetArgs(c, reqUrl)

	// 构造请求
	req := cas.BuildLoginRequest(c, lt, postUrl, reqUrl)

	// 发送请求
	SendRequest(c, req)

	// 读取响应内容
	body := ReadBody(c.Response)

	// 检查标题
	cas.LoginStatusFilterUP(body)

	// 判断本次登陆账号与原账号是否不同
	rid, rsid := gw.IsLoginRepeatedly(body)
	// 若不同
	if len(rid) > 0 {
		// 挤下线
		c.Option.Mute = true
		gw.Replace(c, rid, rsid)
		body = ReadBody(c.Response)
	}

	// 判断是否已欠费
	overdue := gw.IsOverdue(body)
	if overdue {
		// 若已欠费，则本次登陆无效，无需登出，算作登出成功
		// 需要注意输出被登出的账号，而不是拿来完成本机下线操作的账号
		if len(rid) > 0 {
			InfoF(successLogout, rid)
		} else {
			InfoF(successLogout, c.User.Username)
		}
		return
	}

	// 获取现在的SID
	sid, _ := gw.GetSIDAndIP(body)
	// 若获取失败
	if len(sid) < 1 {
		FatalL(failGetInfo)
	}

	// 踢下线
	ok := gw.Kick(c, sid)
	if !ok {
		if len(rid) > 0 {
			FatalF(failLogout, rid)
		}
		FatalF(failLogout, c.User.Username)
	}

	// 请求成功，打印真实的被登出账号
	if len(rid) > 0 {
		InfoF(successLogout, rid)
	} else {
		InfoF(successLogout, c.User.Username)
	}
}

// 使用Cookie登出，基本与使用账号登出一致，但是由于是第二优先级，登出失败时返回bool而不是结束程序
func logoutWithC(c *ctx.Ctx) (ok bool) {
	if ctx.FullView {
		InfoF(usingC, c.User.Cookie.Value)
	}

	// 设置Cookie
	c.Client.Jar.SetCookies(&url.URL{
		Scheme: "https",
		Host:   "ipgw.neu.edu.cn",
	}, []*http.Cookie{c.Net.Cookie})

	// 构造请求
	req, _ := http.NewRequest("GET", "https://ipgw.neu.edu.cn/srun_cas.php?ac_id=1", nil)

	// 发送请求
	SendRequest(c, req)

	// 获取响应体
	body := ReadBody(c.Response)

	// 检查标题
	ok = cas.LoginStatusFilterC(body)
	// 检查未通过
	if !ok {
		return false
	}

	// 判断本次登陆账号与原账号是否不同
	rid, rsid := gw.IsLoginRepeatedly(body)
	// 若不同
	if len(rid) > 0 {
		// 挤下线
		// todo Replace里遇到Kick失败直接就Fatal了，暂时不改逻辑
		c.Option.Mute = true
		gw.Replace(c, rid, rsid)
		body = ReadBody(c.Response)
	}

	// 挂载ID，不需要检查是否获取失败
	c.User.Username = gw.GetID(body)

	// 判断是否已欠费
	overdue := gw.IsOverdue(body)
	if overdue {
		// 若已欠费，则本次登陆无效，无需登出，算作登出成功
		// 需要注意输出被登出的账号，而不是拿来完成本机下线操作的账号
		if len(rid) > 0 {
			InfoF(successLogout, rid)
		} else {
			InfoF(successLogout, c.User.Username)
		}
		return true
	}

	// 获取现在的SID
	sid, _ := gw.GetSIDAndIP(body)
	// 若获取失败
	if len(sid) < 1 {
		ErrorL(failGetInfo)
		return false
	}

	// 踢下线
	ok = gw.Kick(c, sid)
	if !ok {
		if len(rid) > 0 {
			ErrorF(failLogout, rid)
			return false
		}
		ErrorF(failLogout, c.User.Username)
		return false
	}

	// 请求成功，打印真实的被登出账号
	if len(rid) > 0 {
		InfoF(successLogout, rid)
	} else {
		InfoF(successLogout, c.User.Username)
	}
	return true
}
