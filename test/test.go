package test

import (
	"ipgw/base"
	"ipgw/base/cfg"
)

var CmdTest = &base.Command{
	UsageLine: "ipgw test [speed] [-v view all]",
	Short:     "校园网测试",
	Long: `提供对于校园网的测试功能
  -v    输出所有中间信息

  ipgw test
    测试是否连接校园网
  ipgw test -v
    测试是否连接校园网并输出详细中间信息
`,
}

func init() {
	CmdTest.Flag.BoolVar(&cfg.FullView, "v", false, "")

	CmdTest.Run = runTest
}

func runTest(cmd *base.Command, args []string) {
	testImpl()
}
