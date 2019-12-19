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
	"strconv"
	"strings"
)

func loginWithUP(x *ctx.Ctx) {
	client := ctx.GetClient()

	fmt.Printf(usingUP, x.User.Username)

	// 请求获得必要参数
	resp, err := client.Get("https://pass.neu.edu.cn/tpass/login?service=https%3A%2F%2Fipgw.neu.edu.cn%2Fsrun_cas.php%3Fac_id%3D1")
	if err != nil {
		if cfg.FullView {
			fmt.Fprintf(os.Stderr, errWhenReadLT, err)
		}
		fmt.Fprintln(os.Stderr, errNetwork)
		os.Exit(2)
	}

	// 读取响应内容
	body := share.ReadBody(resp)

	// 读取lt
	ltExp := regexp.MustCompile(`name="lt" value="(.+?)"`)
	lt := ltExp.FindAllStringSubmatch(body, -1)[0][1]

	if cfg.FullView {
		fmt.Printf(successGetLT, lt)
	}

	// 拼接data
	data := "rsa=" + x.User.Username + x.User.Password + lt +
		"&ul=" + strconv.Itoa(len(x.User.Username)) +
		"&pl=" + strconv.Itoa(len(x.User.Password)) +
		"&lt=" + lt +
		"&execution=e1s1" +
		"&_eventId=submit"

	// 构造请求
	req, _ := http.NewRequest("POST",
		"https://pass.neu.edu.cn/tpass/login?service=https%3A%2F%2Fipgw.neu.edu.cn%2Fsrun_cas.php%3Fac_id%3D1",
		strings.NewReader(data))

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Host", "pass.neu.edu.cn")
	req.Header.Add("Origin", "https://pass.neu.edu.cn")
	req.Header.Add("Referer", "https://pass.neu.edu.cn/tpass/login?service=https%3A%2F%2Fipgw.neu.edu.cn%2Fsrun_cas.php%3Fac_id%3D3")
	if x.UA != "" {
		req.Header.Add("User-Agent", x.UA)
	}

	if cfg.FullView {
		fmt.Println(sendingRequest)
	}

	// 发送请求
	resp, err = client.Do(req)
	share.ErrWhenReqHandler(err)

	// 读取响应内容
	body = share.ReadBody(resp)

	// 检查标题
	t := share.GetTitle(body)
	if t == "智慧东大--统一身份认证" {
		fmt.Fprintln(os.Stderr, wrongUOrP)
		os.Exit(2)
	}

	if t == "智慧东大" {
		fmt.Fprintln(os.Stderr, wrongLock)
		os.Exit(2)
	}

	if t == "系统提示" {
		fmt.Fprintln(os.Stderr, wrongBan)
		os.Exit(2)
	}

	if strings.Contains(body, "aaa") {
		body = share.CollisionHandler(body)
	}

	out := share.GetIfOut(body)
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

	out := share.GetIfOut(body)
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
