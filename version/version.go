package version

import (
	. "ipgw/base"
)

var CmdVersion = &Command{
	Name:      "version",
	UsageLine: "ipgw version [-l list]",
	Short:     "版本查询",
	Long: `输出ipgw的版本信息
  -l    输出完整版本功能

  ipgw version
    查看版本
  ipgw version -l
    查看当前版本完整功能
`,
}

var l bool

func init() {
	CmdVersion.Flag.BoolVar(&l, "l", false, "")

	CmdVersion.Run = runVersion
}

func runVersion(cmd *Command, args []string) {
	InfoF(`%s
版本: %s [API %s]`, Title, Version, API)

	if l {
		InfoL(detail)
	}
}
