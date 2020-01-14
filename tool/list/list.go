package list

import (
	. "ipgw/base"
)

var CmdToolList = &Command{
	Name:      "list",
	UsageLine: "ipgw tool list [-a all] [-l local]",
	Short:     "查看工具列表",
	Long: `提供查看工具列表的功能
  -a    查看所有工具
  -l    查看本地工具

  ipgw tool list
    查看可用工具(API兼容)
  ipgw tool list -a
    查看所有工具
  ipgw tool list -l
    查看本地已有的工具
`,
}
var (
	a, l bool
)

func init() {
	CmdToolList.Flag.BoolVar(&a, "a", false, "")
	CmdToolList.Flag.BoolVar(&l, "l", false, "")

	CmdToolList.Run = runToolList
}

func runToolList(cmd *Command, args []string) {
	if l {
		// 若l则输出本地
		printLocalList()
	} else {
		// 不然输出在线地址，是否输出全部通过参数传入
		printOnlineList(a)
	}
}
