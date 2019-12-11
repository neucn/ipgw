package version

import (
	"fmt"
	"ipgw/base"
	"ipgw/text"
)

var CmdVersion = &base.Command{
	UsageLine: "ipgw version",
	Short:     "版本查询",
	Long:      "输出ipgw的版本信息, 包括当前版本与草案功能",
}

func init() {
	CmdVersion.Run = runVersion // break init cycle
}

func runVersion(cmd *base.Command, args []string) {
	fmt.Printf("%s\n%s", base.IPGW.Long, text.VersionDetail)
}
