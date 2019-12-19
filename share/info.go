package share

import (
	"fmt"
	"io/ioutil"
	"ipgw/base/cfg"
	"ipgw/base/ctx"
	"math/rand"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func GetIPAndSID(body string, x *ctx.Ctx) (ok bool) {
	// 挂载IP信息
	ipExp := regexp.MustCompile(`get_online_info\('(.+?)'\)`)
	ip := ipExp.FindAllStringSubmatch(body, -1)

	if len(ip) == 0 {
		return false
	}
	x.Net.IP = ip[0][1]
	if cfg.FullView {
		fmt.Printf(successGetIP, x.Net.IP)
	}

	// 挂载SID信息
	// todo 更改匹配方式
	//sidExp := regexp.MustCompile(`do_drop\('(.+?)'\)`)
	sidExp := regexp.MustCompile(`background:lightgreen[\w\W]+?onclick="do_drop\('(\d+)'\)`)
	sidList := sidExp.FindAllStringSubmatch(body, -1)
	if len(sidList) < 1 {
		return false
	}

	x.Net.SID = sidList[0][1]
	if cfg.FullView {
		fmt.Printf(successGetSID, x.Net.SID)
	}

	return true
}

func GetIDAndSIDWhenCollision(body string) (id string, sid string) {
	idExp := regexp.MustCompile(`aaa\n(\d+?)ccc`)
	idList := idExp.FindAllStringSubmatch(body, -1)
	if len(idList) < 1 {
		return "", ""
	}

	id = idList[0][1]

	sidExp := regexp.MustCompile(`btn-dark" href="javascript\(0\);" onclick="do_drop\('(\d+?)'\);`)
	sidList := sidExp.FindAllStringSubmatch(body, -1)
	if len(sidList) < 1 {
		return id, ""
	}

	return id, sidList[0][1]
}

func GetIfOut(body string) (out bool) {
	outExp := regexp.MustCompile(`余额不足月租`)
	return outExp.MatchString(body)
}

func GetTitle(body string) string {
	titleExp := regexp.MustCompile(`<title>(.+?)</title>`)
	title := titleExp.FindAllStringSubmatch(body, -1)
	if len(title) < 1 {
		fmt.Fprintln(os.Stderr, failGetResp)
		os.Exit(2)
	}

	return title[0][1]
}

func ReadBody(resp *http.Response) (body string) {
	res, _ := ioutil.ReadAll(resp.Body)
	_ = resp.Body.Close()
	return string(res)
}

func GetDevice(name string, x *ctx.Ctx) {
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
		fmt.Printf(deviceNotFound, name)
	}

	x.UA = ua
}

func PrintNetInfo(x *ctx.Ctx) {
	if cfg.FullView {
		fmt.Println(gettingInfo)
	}

	// 检查是否登陆
	if x.Net.IP == "" {
		fmt.Fprintln(os.Stderr, errState)
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
	ErrWhenReqHandler(err)
	// 读取响应内容
	body := ReadBody(resp)

	// 解析响应
	split := strings.Split(body, ",")
	if len(split) != 6 {
		fmt.Fprintln(os.Stderr, errState)
		return
	}
	x.Net.Used, err = strconv.Atoi(split[0])
	x.Net.Time, err = strconv.Atoi(split[1])
	x.Net.Balance, err = strconv.ParseFloat(split[2], 64)

	if cfg.FullView {
		fmt.Println(successGetInfo)
	}

	x.Net.Print()
}
