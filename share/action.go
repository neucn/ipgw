package share

import (
	"ipgw/base/ctx"
	"net/http"
	"strings"
)

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
