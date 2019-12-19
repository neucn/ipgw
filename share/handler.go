package share

import (
	"fmt"
	"ipgw/base/cfg"
	"ipgw/base/ctx"
	"os"
)

func ErrWhenReqHandler(err error) {
	if err != nil {
		if cfg.FullView {
			fmt.Fprintf(os.Stderr, errRequest, err)
		}
		fmt.Fprintln(os.Stderr, errNetwork)
		os.Exit(2)
	}
}

func CollisionHandler(body string) string {
	id, sid := GetIDAndSIDWhenCollision(body)
	if id == "" {
		fmt.Fprintln(os.Stderr, errState)
		os.Exit(2)
	}

	fmt.Printf(differentU, id)

	if sid == "" {
		fmt.Fprintln(os.Stderr, failGetInfo)
		os.Exit(2)
	}

	if cfg.FullView {
		fmt.Printf(beginLogout, id)
	}

	client := ctx.GetClient()

	// 踢下线
	resp, err := Kick(sid)

	ErrWhenReqHandler(err)
	body = ReadBody(resp)

	if cfg.FullView {
		fmt.Println(body)
	}

	if body != "下线请求已发送" {
		fmt.Fprintf(os.Stderr, failLogout, id)
		os.Exit(2)
	}

	fmt.Printf(successLogout, id)

	resp, err = client.Get("https://ipgw.neu.edu.cn/srun_cas.php?ac_id=1")
	ErrWhenReqHandler(err)
	return ReadBody(resp)
}
