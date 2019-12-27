package list

import (
	"fmt"
	"ipgw/base/cfg"
	"ipgw/base/ctx"
	"ipgw/share"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func fetchIndexBodyByCAS(x *ctx.Ctx) (ib string) {
	client := ctx.GetClient()
	client.Jar.SetCookies(&url.URL{
		Scheme: "https",
		Host:   "pass.neu.edu.cn",
	}, []*http.Cookie{x.User.CAS})

	if cfg.FullView {
		fmt.Println(fetchingIndexByCAS)
	}
	resp, err := client.Get("http://ipgw.neu.edu.cn:8800/sso/default/neusoft")
	share.ErrWhenReqHandler(err)

	ib = share.ReadBody(resp)
	title := share.GetTitle(ib)
	if title == "智慧东大--统一身份认证" {
		fmt.Println(wrongCASExpired)
		os.Exit(2)
	}

	if cfg.FullView {
		fmt.Println(successFetch)
	}

	return
}

func fetchIndexBodyByUP(x *ctx.Ctx) (ib string) {
	reqUrl := "https://pass.neu.edu.cn/tpass/login?service=http://ipgw.neu.edu.cn:8800/sso/default/neusoft"
	if cfg.FullView {
		fmt.Printf(fetchingIndexByUP, x.User.Username)
	}

	ib = share.Login(reqUrl, x)

	// Login里检查过标题了
	if cfg.FullView {
		fmt.Println(successFetch)
	}
	return
}

func fetchBillBody() (bb string) {
	if cfg.FullView {
		fmt.Println(fetchingBill)
	}
	resp, err := ctx.GetClient().Get("http://ipgw.neu.edu.cn:8800/financial/checkout/list")
	share.ErrWhenReqHandler(err)
	if cfg.FullView {
		fmt.Println(successFetch)
	}
	return share.ReadBody(resp)
}

func fetchRechargeBody() (rb string) {
	if cfg.FullView {
		fmt.Println(fetchingRecharge)
	}
	resp, err := ctx.GetClient().Get("http://ipgw.neu.edu.cn:8800/financial/pay/list")
	share.ErrWhenReqHandler(err)
	if cfg.FullView {
		fmt.Println(successFetch)
	}
	return share.ReadBody(resp)
}

func fetchHistoryBody(p int) (hb string) {
	if cfg.FullView {
		fmt.Println(fetchingRecharge)
	}
	resp, err := ctx.GetClient().Get(fmt.Sprintf("http://ipgw.neu.edu.cn:8800/log/detail/index?page=%d&per-page=20", p))
	share.ErrWhenReqHandler(err)
	if cfg.FullView {
		fmt.Println(successFetch)
	}
	return share.ReadBody(resp)
}

func processUser(body string) {
	// 取出id和name

	fmt.Println("# 基本信息")
	idExp := regexp.MustCompile(`账号</label>\s+(\d+)\s+`)
	nameExp := regexp.MustCompile(`姓名</label>\s+(.+?)\s+`)
	ids := idExp.FindAllStringSubmatch(body, -1)
	names := nameExp.FindAllStringSubmatch(body, -1)
	if len(ids) < 1 || len(names) < 1 {
		fmt.Fprintln(os.Stderr, failToFetch)
		return
	}
	fmt.Printf(`   姓名	%s
   学号	%s
`, names[0][1], ids[0][1])

	fmt.Println()

}

func processDevice(body string) {
	// 取出device

	fmt.Println("# 登陆设备")
	dExp := regexp.MustCompile(`<td>\d+</td>\W+?<td>(.+?)</td>\W+?<td>.+?</td>\W+?<td>(.+?)</td>\W+?<td>(.+?)</td>\W+?<td><a id="(\d+)".+?下线</a></td>`)
	ds := dExp.FindAllStringSubmatch(body, -1)

	if f {
		for i, d := range ds {
			fmt.Printf(`## 设备%d
   IP	%s
   SID	%s
   类型	%s
   时间	%s
`, i, d[1], d[4], d[3], d[2])
		}
	} else {
		for i, d := range ds {
			fmt.Printf("   No.%d\t%s\t%s\t%s\n", i, d[2], d[1], d[4])
		}
	}

	fmt.Println()
}

func processInfo(body string) {
	// 取出套餐信息

	fmt.Println("# 套餐信息")
	infoExp := regexp.MustCompile(`<td>\W+.+?(\d+?)G.+?(\d+?)元.+?</td>\W+<td>\W+(.+?)\W+</td>\W+<td>(.+?)</td>\W+<td>(.+?)</td>\W+<td>.+?</td>\W+<td>(.+?)</td>\W+<td>.+?</td>`)
	infos := infoExp.FindAllStringSubmatch(body, -1)
	if len(infos) < 1 {
		fmt.Fprintln(os.Stderr, failToFetch)
		return
	}

	b, _ := strconv.ParseFloat(infos[0][6], 32)
	var status string
	if b < 0 {
		status = "已欠费"
	} else {
		status = "正常"
	}

	fmt.Printf(`   套餐	%sG / %sR
   已用	%s
   时长	%s
   次数	%s
   余额	%sR
   状态	%s
`, infos[0][1], infos[0][2], infos[0][3], infos[0][4], infos[0][5], infos[0][6], status)

	fmt.Println()
}

func processBill(body string) {
	// 取出付款账单
	fmt.Println("# 扣款记录")
	tExp := regexp.MustCompile(`<title>(.+?)</title>`)
	t := tExp.FindAllStringSubmatch(body, -1)
	if len(t) < 1 || t[0][1] != "结算清单" {
		fmt.Fprintln(os.Stderr, failToFetch)
		return
	}

	billExp := regexp.MustCompile(`<td>(\d+?)</td><td>\d+?</td><td>(\d+?)</td><td>.+?</td><td>.+?</td><td>(.+?)G</td><td>(.+?)</td><td>.+?</td><td>(.+?)</td>`)
	bills := billExp.FindAllStringSubmatch(body, -1)

	if f {
		for _, b := range bills {
			fmt.Printf(`## %s
   扣款	%sR
   流量	%sG
   时长	%s
   流水	%s
`, strings.Split(b[5], " ")[0], b[2], b[3], b[4], b[1])
		}
	} else {
		for _, b := range bills {
			fmt.Printf("   %s\t%sR\t%sG\n", strings.Split(b[5], " ")[0], b[2], b[3])
		}
	}

	fmt.Println()
}

func processHistory(body string) {
	// 取出使用日志
	fmt.Println("# 使用记录")
	tExp := regexp.MustCompile(`<title>(.+?)</title>`)
	t := tExp.FindAllStringSubmatch(body, -1)
	if len(t) < 1 || t[0][1] != "上网明细" {
		fmt.Fprintln(os.Stderr, failToFetch)
		return
	}

	hExp := regexp.MustCompile(`<td>\d+?</td><td>(.+?)</td><td>(.+?)</td><td>(.+?)</td><td>(.+?)</td><td>(.+?)</td><td>.+?</td></tr>`)
	hs := hExp.FindAllStringSubmatch(body, -1)

	for _, h := range hs {
		fmt.Printf("   %s - %s\t%s\t%s\n", h[1], h[2], h[3], h[4])
	}

	fmt.Println()
}

func processRecharge(body string) {
	// 取出充值记录

	fmt.Println("# 充值记录")
	tExp := regexp.MustCompile(`<title>(.+?)</title>`)
	t := tExp.FindAllStringSubmatch(body, -1)
	if len(t) < 1 || t[0][1] != "缴费清单" {
		fmt.Fprintln(os.Stderr, failToFetch)
		return
	}

	rExp := regexp.MustCompile(`<td>(\d+?)</td><td>\d+?</td><td>(\d+?)</td><td>.+?</td><td>(.+?)</td><td>.+?</td><td>(.+?)</td><td>.+?</td>`)
	rs := rExp.FindAllStringSubmatch(body, -1)

	if f {
		for _, r := range rs {
			fmt.Printf(`## 流水 %s
   金额	%sR
   途径	%s
   时间	%s
`, r[1], r[2], r[3], r[4])
		}
	} else {
		for _, r := range rs {
			fmt.Printf("   %s\t%sR\t%s\n", r[4], r[2], r[3])
		}
	}

	fmt.Println()
}

func processLocal(x *ctx.Ctx) {
	fmt.Println("# 本地信息")

	if x.User.Username != "" {
		fmt.Printf("   已保存账号\t%s\n", x.User.Username)
	} else {
		fmt.Printf("   已保存账号\t%s\n", "无")
	}

	if x.User.Cookie.Value != "" {
		fmt.Printf("   已存储Cookie\t%s\n", x.User.Cookie.Value)
	} else {
		fmt.Printf("   已存储Cookie\t%s\n", "无")
	}

	if x.User.CAS.Value != "" {
		fmt.Printf("   已存储CAS\t%s\n", x.User.CAS.Value)
	} else {
		fmt.Printf("   已存储CAS\t%s\n", "无")
	}

	fmt.Println()
}
