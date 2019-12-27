package version

import (
	"fmt"
	"ipgw/base"
)

var CmdVersion = &base.Command{
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

var u, l bool

func init() {
	CmdVersion.Flag.BoolVar(&l, "l", false, "")

	CmdVersion.Run = runVersion // break init cycle
}

func runVersion(cmd *base.Command, args []string) {
	fmt.Println(base.IPGW.Long)

	if l {
		fmt.Println(detail)
	}
}
