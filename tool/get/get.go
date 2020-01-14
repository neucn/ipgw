package get

import (
	"flag"
	. "ipgw/base"
	"ipgw/ctx"
	"ipgw/help"
)

var CmdToolGet = &Command{
	Name:      "get",
	UsageLine: "ipgw tool get tool1 tool2 ...",
	Short:     "下载工具",
	Long: `提供下载工具的功能

  ipgw tool get teemo
    下载teemo
`,
}

func init() {
	flag.BoolVar(&ctx.FullView, "v", false, "")

	CmdToolGet.Run = runToolGet
}

func runToolGet(cmd *Command, args []string) {
	// 如果没有指定工具名
	if len(args) < 1 {
		help.PrintSimpleUsage(cmd)
	}

	// 处理业务逻辑
	getAll(args)
}
