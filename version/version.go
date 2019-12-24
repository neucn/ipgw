package version

import (
	"fmt"
	"io/ioutil"
	"ipgw/base"
	"ipgw/base/cfg"
	"ipgw/base/ctx"
	"os"
	"regexp"
)

var CmdVersion = &base.Command{
	UsageLine: "ipgw version",
	Short:     "版本查询",
	Long: `输出ipgw的版本信息
  -u    查看最新版本
  -v    输出完整版本功能

  ipgw version
    查看版本
  ipgw version -u
    查看最新版本
  ipgw version -v
    查看当前版本完整功能
`,
}

var u bool

func init() {
	CmdVersion.Flag.BoolVar(&u, "u", false, "")
	CmdVersion.Flag.BoolVar(&cfg.FullView, "v", false, "")

	CmdVersion.Run = runVersion // break init cycle
}

func runVersion(cmd *base.Command, args []string) {
	fmt.Println(base.IPGW.Long)

	if cfg.FullView {
		fmt.Println(detail)
	}

	if u {
		client := ctx.GetClient()
		fmt.Println(tipQuery)

		resp, err := client.Get("https://api.github.com/repos/imyown/ipgw/releases/latest")
		if err != nil {
			fmt.Fprintln(os.Stderr, errNet)
			os.Exit(2)
		}

		res, err := ioutil.ReadAll(resp.Body)
		_ = resp.Body.Close()
		body := string(res)

		tagExp := regexp.MustCompile(`"tag_name":"(.+?)"`)
		tags := tagExp.FindAllStringSubmatch(body, -1)

		if len(tags) == 0 {
			fmt.Println(failQuery)
			return
		}

		tag := tags[0][1]
		if tag == cfg.Version {
			fmt.Println(tipAlreadyLatest)
			return
		}
		fmt.Printf(tipLatest, tag)
	}

}
