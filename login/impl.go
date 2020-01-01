package login

import (
	"fmt"
	"ipgw/base/cfg"
	"ipgw/base/ctx"
	"ipgw/share"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
)

func loginWithUP(x *ctx.Ctx) {
	client := ctx.GetClient()
	reqUrl := "https://pass.neu.edu.cn/tpass/login?service=https%3A%2F%2Fipgw.neu.edu.cn%2Fsrun_cas.php%3Fac_id%3D1"

	fmt.Printf(usingUP, x.User.Username)

	body := share.Login(reqUrl, x)

	// 检查内容
	if strings.Contains(body, "aaa") {
		body = share.CollisionHandler(body)
	}

	out := share.GetIfUsedOut(body)
	if out {
		fmt.Println(failBalanceOut)
		os.Exit(2)
	}

	// 读取IP与SID
	ok := share.GetIPAndSID(body, x)
	if !ok {
		fmt.Fprintln(os.Stderr, failLogin)
		os.Exit(2)
	}

	cookie := client.Jar.Cookies(&url.URL{
		Scheme: "https",
		Host:   "ipgw.neu.edu.cn",
	})

	cases := client.Jar.Cookies(&url.URL{
		Scheme: "https",
		Host:   "pass.neu.edu.cn",
		Path:   "/tpass/",
	})

	for _, cas := range cases {
		if cas.Name == "CASTGC" {
			x.User.CAS = cas
			break
		}
	}

	if len(cookie) == 0 {
		fmt.Fprintln(os.Stderr, failGetCookie)
	} else {
		x.User.Cookie = cookie[0]
		if cfg.FullView {
			fmt.Printf(successGetCookie, x.User.Cookie.Value)
		}
	}

	fmt.Printf(successLogin, x.User.Username)
}

func loginWithC(x *ctx.Ctx) {
	client := ctx.GetClient()

	fmt.Printf(usingC, x.User.Cookie.Value)

	// 请求获得必要参数
	client.Jar.SetCookies(&url.URL{
		Scheme: "https",
		Host:   "ipgw.neu.edu.cn",
	}, []*http.Cookie{x.User.Cookie})

	if cfg.FullView {
		fmt.Println(sendingRequest)
	}

	var resp *http.Response
	var err error

	if x.UA == "" {
		resp, err = client.Get("https://ipgw.neu.edu.cn/srun_cas.php?ac_id=1")
	} else {
		req, _ := http.NewRequest("GET", "https://ipgw.neu.edu.cn/srun_cas.php?ac_id=1", nil)
		req.Header.Add("User-Agent", x.UA)

		resp, err = client.Do(req)
	}

	share.ErrWhenReqHandler(err)

	// 读取响应内容
	body := share.ReadBody(resp)

	// 检查标题
	t := share.GetTitle(body)
	if t == "智慧东大--统一身份认证" {
		fmt.Fprintln(os.Stderr, failCookieExpired)
		os.Exit(2)
	}

	if strings.Contains(body, "aaa") {
		body = share.CollisionHandler(body)
	}

	// 读取学号
	usernameExp := regexp.MustCompile(`user_name" style="float:right;color: #894324;">(.+?)</span>`)
	username := usernameExp.FindAllStringSubmatch(body, -1)

	if len(username) == 0 {
		fmt.Println(failGetInfo)
	} else {
		x.User.Username = username[0][1]
		if cfg.FullView {
			fmt.Printf(successGetUsername, x.User.Username)
		}
	}

	out := share.GetIfUsedOut(body)
	if out {
		fmt.Println(failBalanceOut)
		os.Exit(2)
	}

	// 读取IP与SID
	ok := share.GetIPAndSID(body, x)

	if !ok {
		fmt.Fprintln(os.Stderr, failLogin)
		os.Exit(2)
	}

	fmt.Printf(successLogin, x.User.Username)

}
