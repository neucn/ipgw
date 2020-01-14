package fix

import . "ipgw/base"

var CmdFix = &Command{
	Name:      "fix",
	UsageLine: "ipgw fix",
	Short:     "修复配置文件",
	Long: `使用空配置覆盖配置文件

  ipgw fix
    修复配置文件
`,
}

func init() {
	CmdFix.Run = runFix
}

func runFix(cmd *Command, args []string) {
	fix()
}
