package update

import (
	"fmt"
	"ipgw/base"
	"ipgw/base/cfg"
)

var CmdUpdate = &base.Command{
	UsageLine: "ipgw update [-f force] [-v view all]",
	Short:     "更新版本",
	Long: `将ipgw更新到最新版
  -f    强制更新
  -v    输出中间信息与详细报错信息

  ipgw update
    检查更新并自动更新
  ipgw update -f
    强制下载最新版更新
`,
}

var f bool

func init() {
	CmdUpdate.Flag.BoolVar(&f, "f", false, "")
	CmdUpdate.Flag.BoolVar(&cfg.FullView, "v", false, "")

	CmdUpdate.Run = runUpdate // break init cycle
}

func runUpdate(cmd *base.Command, args []string) {
	fmt.Printf(localVersion, cfg.Version)

	c := checkVersion()
	// todo  更新自己
	if f || c.Update {
		if f {
			fmt.Println(forcing)
		} else {
			printChangelog(c)
		}
		update(c)
	}
}
