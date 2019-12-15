package kick

import (
	"fmt"
	"io/ioutil"
	"ipgw/base/cfg"
	"ipgw/base/share"
	"os"
)

func kickWithSID(sid string) {
	fmt.Printf(tipBeginKick, sid)

	resp, err := share.Kick(sid)

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
