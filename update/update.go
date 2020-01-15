package update

import (
	. "ipgw/base"
	"ipgw/ctx"
	. "ipgw/lib"
)

var CmdUpdate = &Command{
	Name:      "update",
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
	CmdUpdate.Flag.BoolVar(&ctx.FullView, "v", false, "")

	CmdUpdate.Run = runUpdate
}

func runUpdate(cmd *Command, args []string) {
	InfoF(localVersion, Version)

	// 获取版本信息
	v := checkVersion()

	// 如果有更新或者强制更新
	if f || v.Update {
		if f {
			InfoL(forcing)
		} else {
			PrintChangelog(Version, v)
		}
		update(v)
	}
}
