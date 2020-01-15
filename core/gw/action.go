package gw

import (
	. "ipgw/base"
	. "ipgw/core/global"
	"ipgw/ctx"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
)

// 强制指定sid下线
func Kick(c *ctx.Ctx, sid string) (ok bool) {
	// 构造请求
	url := "https://ipgw.neu.edu.cn/srun_cas.php"
	data := "action=dm&sid=" + sid
	req, _ := http.NewRequest("POST", url, strings.NewReader(data))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Host", "ipgw.neu.edu.cn")
	req.Header.Set("Origin", " https://ipgw.neu.edu.cn")
	req.Header.Set("Referer", "https://ipgw.neu.edu.cn/srun_cas.php?ac_id=1")

	// 发送请求
	SendRequest(c, req)

	// 响应体
	result := ReadBody(c.Response)
	if !c.Option.Mute && ctx.FullView {
		// 输出响应信息
		InfoL(result)
	}

	// 判断是否成功
	if result == "下线请求已发送" {
		return true
	}
	return false
}

// 解决重复登陆，先强制指定sid下线再重新访问页面获取信息，更新ctx的Response
func Replace(c *ctx.Ctx, id, sid string) {
	mute := c.Option.Mute
	if !mute && ctx.FullView {
		InfoF(differentU, id)
	}

	// 强制下线
	ok := Kick(c, sid)
	if !ok {
		// 若强制下线失败，结束程序
		FatalF(failLogout, id)
	}

	if !mute {
		InfoF(successLogout, id)
	}

	// 重新访问页面
	req, _ := http.NewRequest("GET", "https://ipgw.neu.edu.cn/srun_cas.php?ac_id=1", nil)
	SendRequest(c, req)
}

// 获取并挂载网络信息
func GetNetInfo(c *ctx.Ctx) {
	mute := c.Option.Mute
	if !mute && ctx.FullView {
		InfoL(infoFetchingNetInfo)
	}

	// 构造请求
	k := strconv.Itoa(rand.Intn(100000 + 1))
	data := "action=get_online_info&key=" + k

	req, _ := http.NewRequest("POST", "https://ipgw.neu.edu.cn/include/auth_action.php?k="+k, strings.NewReader(data))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Add("Host", "pass.neu.edu.cn")
	req.Header.Add("Origin", "https://pass.neu.edu.cn")
	req.Header.Add("Referer", "https://ipgw.neu.edu.cn/srun_cas.php?ac_id=1")

	// 发送请求
	SendRequest(c, req)

	// 读取响应内容
	body := ReadBody(c.Response)

	// 解析响应
	split := strings.Split(body, ",")

	// 可能有登陆异常的情况
	if len(split) != 6 {
		ErrorL(errState)
		return
	}

	// 挂载
	c.Net.Used, _ = strconv.Atoi(split[0])
	c.Net.Time, _ = strconv.Atoi(split[1])
	c.Net.Balance, _ = strconv.ParseFloat(split[2], 64)

	if !mute && ctx.FullView {
		InfoL(successGetInfo)
	}
}
