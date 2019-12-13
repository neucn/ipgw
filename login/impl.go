package login

import (
	"fmt"
	"io/ioutil"
	"ipgw/base/cfg"
	"ipgw/base/ctx"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func loginWithUP(x *ctx.Ctx) {
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
		fmt.Println(tipRequest)
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
	getSIDAndIP(body, x)

	cookie := client.Jar.Cookies(&url.URL{
		Scheme: "https",
		Host:   "ipgw.neu.edu.cn",
	})

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

	fmt.Printf(tipBeginWithC, x.User.Cookie.Value)

	// 请求获得必要参数
	client.Jar.SetCookies(&url.URL{
		Scheme: "https",
		Host:   "ipgw.neu.edu.cn",
	}, []*http.Cookie{x.User.Cookie})

	if cfg.FullView {
		fmt.Println(tipRequest)
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

	if err != nil {
		fmt.Printf(errWhenRequest, err)
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
		os.Exit(2)
	} else {
		x.User.Username = username[0][1]
		if cfg.FullView {
			fmt.Printf(successGetUsername, x.User.Username)
		}
	}

	// 读取IP与SID
	getSIDAndIP(body, x)

	fmt.Printf(successLogin, x.User.Username)

}

func getSIDAndIP(body string, x *ctx.Ctx) {
	// 挂载IP信息
	ipExp := regexp.MustCompile(`get_online_info\('(.+?)'\)`)
	ip := ipExp.FindAllStringSubmatch(body, -1)

	if len(ip) == 0 {
		fmt.Fprintln(os.Stderr, wrongUOrP)
		os.Exit(2)
	} else {
		x.Net.IP = ip[0][1]
		if cfg.FullView {
			fmt.Printf(successGetIP, x.Net.IP)
		}
	}

	// 挂载SID信息
	// todo 更改匹配方式
	sidExp := regexp.MustCompile(`do_drop\('(.+?)'\)`)
	sidList := sidExp.FindAllStringSubmatch(body, -1)
	sid := sidList[len(sidList)-1][1]
	if sid == "" {
		fmt.Fprintln(os.Stderr, failGetSID)
	} else {
		x.Net.SID = sid
		if cfg.FullView {
			fmt.Printf(successGetSID, x.Net.SID)
		}
	}
}

func getDevice(name string, x *ctx.Ctx) {
	if name == "" {
		return
	}

	var ua string

	// todo 由于移动端会跳转到单独网页，因此暂时把移动端去掉
	switch name {
	case "win":
		fallthrough
	case "windows":
		ua = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/79.0.3945.79 Safari/537.36"
	case "linux":
		ua = "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/74.0.3729.131 Safari/537.36"
	case "osx":
		fallthrough
	case "darwin":
		ua = "Opera/9.80 (Macintosh; Intel Mac OS X 10.6.8; U; en) Presto/2.8.131 Version/11.11"
	//case "ios":
	//	ua = "Mozilla/5.0 (iPhone; CPU iPhone OS 11_0_2 like Mac OS X) AppleWebKit/604.1.38 (KHTML, like Gecko) Mobile/15A421 wxwork/2.5.8 MicroMessenger/6.3.22 Language/zh"
	//case "android":
	//	ua = "Mozilla/5.0 (Linux; Android 8.0; DUK-AL20 Build/HUAWEIDUK-AL20; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/57.0.2987.132 MQQBrowser/6.2 TBS/044353 Mobile Safari/537.36 MicroMessenger/6.7.3.1360(0x26070333) NetType/WIFI Language/zh_CN Process/tools"
	default:
		fmt.Printf("暂不支持 %s, 自动转换为匿名设备\n", name)
	}

	x.UA = ua
}

func printNetInfo(x *ctx.Ctx) {
	if cfg.FullView {
		fmt.Println("获取账号信息中...")
	}

	// 检查是否登陆
	if x.Net.IP == "" {
		fmt.Println("登陆状态异常")
		return
	}

	// 获取client实例
	client := ctx.GetClient()

	// 构造请求
	k := strconv.Itoa(rand.Intn(100000 + 1))
	data := "action=get_online_info&key=" + k

	req, _ := http.NewRequest("POST", "https://ipgw.neu.edu.cn/include/auth_action.php?k="+k, strings.NewReader(data))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Add("Host", "pass.neu.edu.cn")
	req.Header.Add("Origin", "https://pass.neu.edu.cn")
	req.Header.Add("Referer", "https://ipgw.neu.edu.cn/srun_cas.php?ac_id=1")

	// 发送请求
	resp, err := client.Do(req)
	if err != nil {
		if cfg.FullView {
			fmt.Fprintf(os.Stderr, "遇到异常: %s\n", err)
			return
		}
		fmt.Println("请检查网络连接")
		return
	}

	// 读取响应内容
	res, err := ioutil.ReadAll(resp.Body)
	_ = resp.Body.Close()
	body := string(res)

	// 解析响应
	split := strings.Split(body, ",")
	if len(split) != 6 {
		fmt.Println("登陆状态异常")
		return
	}
	x.Net.Used, err = strconv.Atoi(split[0])
	x.Net.Time, err = strconv.Atoi(split[1])
	x.Net.Balance, err = strconv.ParseFloat(split[2], 64)

	if cfg.FullView {
		fmt.Println("获取信息成功")
	}

	x.Net.Print()
}
