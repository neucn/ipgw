package share

import (
	"fmt"
	"ipgw/base/cfg"
	"ipgw/base/ctx"
	"net/http"
	"strconv"
	"strings"
)

func Login(reqUrl string, x *ctx.Ctx) (body string) {
	client := ctx.GetClient()
	// 请求获得必要参数
	lt, postUrl := GetLTAndURL(reqUrl)
	// 拼接data
	data := "rsa=" + x.User.Username + x.User.Password + lt +
		"&ul=" + strconv.Itoa(len(x.User.Username)) +
		"&pl=" + strconv.Itoa(len(x.User.Password)) +
		"&lt=" + lt +
		"&execution=e1s1" +
		"&_eventId=submit"

	// 构造请求
	req, _ := http.NewRequest("POST",
		"https://pass.neu.edu.cn"+postUrl,
		strings.NewReader(data))

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Host", "pass.neu.edu.cn")
	req.Header.Add("Origin", "https://pass.neu.edu.cn")
	req.Header.Add("Referer", reqUrl)
	if x.UA != "" {
		req.Header.Add("User-Agent", x.UA)
	}

	if cfg.FullView {
		fmt.Println(sendingRequest)
	}

	// 发送请求
	resp, err := client.Do(req)
	ErrWhenReqHandler(err)

	// 读取响应内容
	body = ReadBody(resp)

	// 检查标题
	t := GetTitle(body)
	TitleAfterLoginHandler(t)

	return
}

func Kick(sid string) (resp *http.Response, err error) {
	// 获取client
	client := ctx.GetClient()

	// 构造请求
	url := "https://ipgw.neu.edu.cn/srun_cas.php"
	data := "action=dm&sid=" + sid
	req, _ := http.NewRequest("POST", url, strings.NewReader(data))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Host", "ipgw.neu.edu.cn")
	req.Header.Set("Origin", " https://ipgw.neu.edu.cn")
	req.Header.Set("Referer", "https://ipgw.neu.edu.cn/srun_cas.php?ac_id=1")

	// 发送请求
	resp, err = client.Do(req)
	return
}
