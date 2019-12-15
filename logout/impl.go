package logout

import (
	"fmt"
	"io/ioutil"
	"ipgw/base/cfg"
	"ipgw/base/ctx"
	"ipgw/base/share"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func logoutWithUP(x *ctx.Ctx) {
	client := ctx.GetClient()

	fmt.Printf(tipBeginWithUP, x.User.Username)

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
	res, err := ioutil.ReadAll(resp.Body)
	_ = resp.Body.Close()
	body := string(res)

	// 读取lt post_url
	ltExp := regexp.MustCompile(`name="lt" value="(.+?)"`)
	lt := ltExp.FindAllStringSubmatch(body, -1)[0][1]

	postUrlExp := regexp.MustCompile(`id="loginForm" action="(.+?)"`)
	postUrl := postUrlExp.FindAllStringSubmatch(body, -1)[0][1]

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
	req, _ := http.NewRequest("POST", "https://pass.neu.edu.cn"+postUrl, strings.NewReader(data))

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Host", "pass.neu.edu.cn")
	req.Header.Add("Origin", "https://pass.neu.edu.cn")
	req.Header.Add("Referer", "https://pass.neu.edu.cn/tpass/login?service=https%3A%2F%2Fipgw.neu.edu.cn%2Fsrun_cas.php%3Fac_id%3D3")
	if x.UA != "" {
		req.Header.Add("User-Agent", x.UA)
	}

	if cfg.FullView {
		fmt.Println(tipGetSID)
	}

	// 发送请求
	resp, err = client.Do(req)

	if err != nil {
		if cfg.FullView {
			fmt.Fprintf(os.Stderr, errWhenRequest, err)
		}
		fmt.Fprintln(os.Stderr, errNetwork)
		os.Exit(2)
	}

	// 读取响应内容
	res, err = ioutil.ReadAll(resp.Body)
	_ = resp.Body.Close()
	body = string(res)

	// 读取IP与SID
	ok := share.GetSIDAndIP(body, x)
	if !ok {
		os.Exit(2)
	}

	resp, err = share.Kick(x.Net.SID)
	if err != nil {
		if cfg.FullView {
			fmt.Fprintf(os.Stderr, errWhenRequest, err)
		}
		fmt.Fprintln(os.Stderr, errNetwork)
		os.Exit(2)
	}

	res, err = ioutil.ReadAll(resp.Body)
	_ = resp.Body.Close()
	body = string(res)

	if cfg.FullView {
		fmt.Println(body)
	}

	if body != "下线请求已发送" {
		fmt.Fprintf(os.Stderr, failLogout, x.Net.SID)
		os.Exit(2)
	}

	fmt.Printf(successLogout, x.User.Username)
}

func logoutWithC(x *ctx.Ctx) (ok bool) {
	client := ctx.GetClient()

	if cfg.FullView {
		fmt.Printf("正在使用Cookie登出 %s\n", x.User.Cookie.Value)
		fmt.Println("获取必要参数中...")
	} else {
		fmt.Println("正在尝试使用Cookie登出...")
	}

	// 请求获得必要参数
	client.Jar.SetCookies(&url.URL{
		Scheme: "https",
		Host:   "ipgw.neu.edu.cn",
	}, []*http.Cookie{x.User.Cookie})

	resp, err := client.Get("https://ipgw.neu.edu.cn/srun_cas.php?ac_id=1")

	if err != nil {
		fmt.Fprintf(os.Stderr, "发送请求时错误: %v\n", err)
		os.Exit(2)
	}

	// 读取响应内容
	res, err := ioutil.ReadAll(resp.Body)
	_ = resp.Body.Close()
	body := string(res)

	// 读取学号
	usernameExp := regexp.MustCompile(`user_name" style="float:right;color: #894324;">(.+?)</span>`)
	username := usernameExp.FindAllStringSubmatch(body, -1)

	if len(username) == 0 {
		fmt.Fprintln(os.Stderr, failCookieExpired)
		return false
	} else {
		x.User.Username = username[0][1]
		if cfg.FullView {
			fmt.Printf("成功获得学号: %s\n", x.User.Username)
		}
	}

	share.GetSIDAndIP(body, x)

	if cfg.FullView {
		fmt.Println("发送登出请求中...")
	}

	resp, err = share.Kick(x.Net.SID)

	if err != nil {
		if cfg.FullView {
			fmt.Fprintf(os.Stderr, errWhenRequest, err)
		}
		fmt.Fprintln(os.Stderr, errNetwork)
		return
	}

	res, err = ioutil.ReadAll(resp.Body)
	_ = resp.Body.Close()
	body = string(res)

	if cfg.FullView {
		fmt.Println(body)
	}

	if body != "下线请求已发送" {
		fmt.Fprintf(os.Stderr, failLogout, x.Net.SID)
		return false
	}

	fmt.Printf(successLogout, x.User.Username)
	return true
}

// Deprecated
// 无法判断是否正确退出
func _logoutWithC(x *ctx.Ctx) {
	client := ctx.GetClient()

	fmt.Printf("正在使用Cookie登出 %s\n", x.User.Cookie.Value)

	// 请求获得必要参数
	client.Jar.SetCookies(&url.URL{
		Scheme: "https",
		Host:   "ipgw.neu.edu.cn",
	}, []*http.Cookie{x.User.Cookie})

	if cfg.FullView {
		fmt.Println("发送登出请求中...")
	}

	resp, err := client.Get("https://ipgw.neu.edu.cn/srun_cas.php?logout")

	if err != nil {
		fmt.Fprintf(os.Stderr, "登出时遇到意外错误: %s", err)
	}

	if resp.StatusCode == http.StatusOK {
		fmt.Printf("登出成功")
	}
}
