package kick

import (
	"fmt"
	"io/ioutil"
	"ipgw/base/cfg"
	"ipgw/base/ctx"
	"net/http"
	"os"
	"strings"
)

func kickWithSID(sid string) {
	fmt.Printf(tipBeginKick, sid)
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
	resp, err := client.Do(req)
	if err != nil {
		if cfg.FullView {
			fmt.Fprintf(os.Stderr, errWhenKick, err)
		}
		fmt.Fprintln(os.Stderr, tipCheckNet)
		return
	}

	res, err := ioutil.ReadAll(resp.Body)
	_ = resp.Body.Close()
	body := string(res)

	if cfg.FullView {
		fmt.Println(body)
	}

	if body != "下线请求已发送" {
		fmt.Fprintf(os.Stderr, failKick, sid)
		return
	}

	fmt.Printf(successKick, sid)
}
