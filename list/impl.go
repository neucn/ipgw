package list

import (
	"fmt"
	. "ipgw/base"
	"ipgw/core/cas"
	. "ipgw/core/global"
	"ipgw/ctx"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
)

// 使用Cookie获取首页信息
func fetchIndexBodyByC(c *ctx.Ctx) (ib string) {
	client := c.Client

	// 设置Cookie
	client.Jar.SetCookies(&url.URL{
		Scheme: "https",
		Host:   "pass.neu.edu.cn",
	}, []*http.Cookie{c.User.Cookie})

	if ctx.FullView {
		InfoL(fetchingIndexByCAS)
	}

	// 构造请求
	req, _ := http.NewRequest("GET", "http://ipgw.neu.edu.cn:8800/sso/default/neusoft", nil)

	// 发送请求
	// 服务器自动跳转到一网通鉴权，因此可以直接访问
	SendRequest(c, req)
	ib = ReadBody(c.Response)

	// 检查是否登陆成功
	cas.LoginStatusFilterC(ib)

	if ctx.FullView {
		InfoL(successFetch)
	}

	return
}

// 使用账号获取首页信息
func fetchIndexBodyByUP(c *ctx.Ctx) (ib string) {
	reqUrl := "https://pass.neu.edu.cn/tpass/login?service=http://ipgw.neu.edu.cn:8800/sso/default/neusoft"
	if ctx.FullView {
		InfoL(fetchingIndexByUP, c.User.Username)
	}

	// 登陆流程
	// 获取必要参数
	lt, postUrl := cas.GetArgs(c, reqUrl)

	// 产生请求
	req := cas.BuildLoginRequest(c, lt, postUrl, reqUrl)

	// 发送请求
	cas.LoginCAS(c, req)

	// 读取响应体
	ib = ReadBody(c.Response)

	// 判断登陆状态
	cas.LoginStatusFilterUP(ib)

	if ctx.FullView {
		InfoL(successFetch)
	}

	return
}

// 获取扣款记录
func fetchBillBody(c *ctx.Ctx) (bb string) {
	if ctx.FullView {
		InfoL(fetchingBill)
	}
	req, _ := http.NewRequest("GET", "http://ipgw.neu.edu.cn:8800/financial/checkout/list", nil)

	SendRequest(c, req)

	if ctx.FullView {
		InfoL(successFetch)
	}
	return ReadBody(c.Response)
}

// 获取充值记录
func fetchRechargeBody(c *ctx.Ctx) (rb string) {
	if ctx.FullView {
		InfoL(fetchingRecharge)
	}
	req, _ := http.NewRequest("GET", "http://ipgw.neu.edu.cn:8800/financial/pay/list", nil)
	SendRequest(c, req)

	if ctx.FullView {
		InfoL(successFetch)
	}
	return ReadBody(c.Response)
}

// 获取使用日志
func fetchHistoryBody(c *ctx.Ctx, p int) (hb string) {
	if ctx.FullView {
		InfoL(fetchingHistory)
	}

	req, _ := http.NewRequest("GET", fmt.Sprintf("http://ipgw.neu.edu.cn:8800/log/detail/index?page=%d&per-page=20", p), nil)

	SendRequest(c, req)

	if ctx.FullView {
		InfoL(successFetch)
	}
	return ReadBody(c.Response)
}

// 打印id和name
func processUser(body string) {

	InfoL("# 基本信息")
	idExp := regexp.MustCompile(`账号</label>\s+(\d+)\s+`)
	nameExp := regexp.MustCompile(`姓名</label>\s+(.+?)\s+`)
	ids := idExp.FindAllStringSubmatch(body, -1)
	names := nameExp.FindAllStringSubmatch(body, -1)
	if len(ids) < 1 || len(names) < 1 {
		ErrorL(failToFetch)
		return
	}
	InfoF(`   姓名	%s
   学号	%s
`, names[0][1], ids[0][1])

	InfoL()

}

// 打印设备信息
func processDevice(body string) {
	// 取出device

	InfoL("# 登陆设备")
	dExp := regexp.MustCompile(`<td>\d+</td>\W+?<td>(.+?)</td>\W+?<td>.+?</td>\W+?<td>(.+?)</td>\W+?<td>(.+?)</td>\W+?<td><a id="(\d+)".+?下线</a></td>`)
	ds := dExp.FindAllStringSubmatch(body, -1)

	if f {
		for i, d := range ds {
			InfoF(`## 设备%d
   IP	%s
   SID	%s
   类型	%s
   时间	%s
`, i, d[1], d[4], d[3], d[2])
		}
	} else {
		for i, d := range ds {
			InfoF("   No.%d\t%s\t%s\t%s\n", i, d[2], d[1], d[4])
		}
	}

	InfoL()
}

// 打印套餐信息
func processInfo(body string) {

	InfoL("# 套餐信息")
	infoExp := regexp.MustCompile(`<td>\W+.+?(\d+?)G下行流量(.+?)元?/.+?</td>\W+<td>\W+(.+?)\W+</td>\W+<td>(.+?)</td>\W+<td>(.+?)</td>\W+<td>.+?</td>\W+<td>(.+?)</td>\W+<td>.+?</td>`)
	infos := infoExp.FindAllStringSubmatch(body, -1)
	if len(infos) < 1 {
		ErrorL(failToFetch)
		return
	}

	b, _ := strconv.ParseFloat(infos[0][6], 32)
	var status string
	if b < 0 {
		status = "已欠费"
	} else {
		status = "正常"
	}

	if infos[0][2] == "免费" {
		infos[0][2] = "0"
	}

	if strings.HasSuffix(infos[0][3], "G") {
		u, _ := strconv.ParseFloat(strings.TrimSuffix(infos[0][6], "G"), 32)
		t, _ := strconv.ParseFloat(infos[0][1], 32)
		if u > t {
			status += "【流量超额】"
		}
	}

	InfoF(`   套餐	%sG / %sR
   已用	%s
   时长	%s
   次数	%s
   余额	%sR
   状态	%s
`, infos[0][1], infos[0][2], infos[0][3], infos[0][4], infos[0][5], infos[0][6], status)

	InfoL()
}

// 打印付款账单
func processBill(body string) {
	InfoL("# 扣款记录")
	tExp := regexp.MustCompile(`<title>(.+?)</title>`)
	t := tExp.FindAllStringSubmatch(body, -1)
	if len(t) < 1 || t[0][1] != "结算清单" {
		ErrorL(failToFetch)
		return
	}

	billExp := regexp.MustCompile(`<td>(\d+?)</td><td>\d+?</td><td>(\d+?)</td><td>.+?</td><td>.+?</td><td>(.+?)G</td><td>(.+?)</td><td>.+?</td><td>(.+?)</td>`)
	bills := billExp.FindAllStringSubmatch(body, -1)

	if f {
		for _, b := range bills {
			InfoF(`## %s
   扣款	%sR
   流量	%sG
   时长	%s
   流水	%s
`, strings.Split(b[5], " ")[0], b[2], b[3], b[4], b[1])
		}
	} else {
		for _, b := range bills {
			InfoF("   %s\t%sR\t%sG\n", strings.Split(b[5], " ")[0], b[2], b[3])
		}
	}

	InfoL()
}

// 打印使用日志
func processHistory(body string) {
	InfoL("# 使用记录")
	tExp := regexp.MustCompile(`<title>(.+?)</title>`)
	t := tExp.FindAllStringSubmatch(body, -1)
	if len(t) < 1 || t[0][1] != "上网明细" {
		ErrorL(failToFetch)
		return
	}

	hExp := regexp.MustCompile(`<td>\d+?</td><td>(.+?)</td><td>(.+?)</td><td>(.+?)</td><td>(.+?)</td><td>(.+?)</td><td>.+?</td></tr>`)
	hs := hExp.FindAllStringSubmatch(body, -1)

	for _, h := range hs {
		InfoF("   %s - %s\t%s\t%s\n", h[1], h[2], h[3], h[4])
	}

	InfoL()
}

// 打印充值记录
func processRecharge(body string) {
	InfoL("# 充值记录")
	tExp := regexp.MustCompile(`<title>(.+?)</title>`)
	t := tExp.FindAllStringSubmatch(body, -1)
	if len(t) < 1 || t[0][1] != "缴费清单" {
		ErrorL(failToFetch)
		return
	}

	rExp := regexp.MustCompile(`<td>(\d+?)</td><td>\d+?</td><td>(\d+?)</td><td>.+?</td><td>(.+?)</td><td>.+?</td><td>(.+?)</td><td>.+?</td>`)
	rs := rExp.FindAllStringSubmatch(body, -1)

	if f {
		for _, r := range rs {
			InfoF(`## 流水 %s
   金额	%sR
   途径	%s
   时间	%s
`, r[1], r[2], r[3], r[4])
		}
	} else {
		for _, r := range rs {
			InfoF("   %s\t%sR\t%s\n", r[4], r[2], r[3])
		}
	}

	InfoL()
}

// 打印本地信息
func processLocal(x *ctx.Ctx) {
	InfoL("# 本地信息")

	if x.User.Username != "" {
		InfoF("   已保存账号\t%s\n", x.User.Username)
	} else {
		InfoF("   已保存账号\t%s\n", "无")
	}

	if x.User.Cookie.Value != "" {
		InfoF("   一网通Cookie\t%s\n", x.User.Cookie.Value)
	} else {
		InfoF("   一网通Cookie\t%s\n", "无")
	}

	if x.Net.Cookie.Value != "" {
		InfoF("   网关Cookie\t%s\n", x.Net.Cookie.Value)
	} else {
		InfoF("   网关Cookie\t%s\n", "无")
	}

	InfoL()
}
