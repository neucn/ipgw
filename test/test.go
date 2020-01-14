package test

import (
	. "ipgw/base"
	"ipgw/ctx"
)

var CmdTest = &Command{
	Name:      "test",
	UsageLine: "ipgw test [-v view all]",
	Short:     "校园网测试",
	Long: `提供对于校园网的测试功能
  -v    输出所有中间信息

  ipgw test
    测试校园网连接与登陆情况
  ipgw test -v
    测试校园网连接与登陆情况并输出详细中间信息
`,
}

func init() {
	CmdTest.Flag.BoolVar(&ctx.FullView, "v", false, "")

	CmdTest.Run = runTest
}

func runTest(cmd *Command, args []string) {
	test()
}
