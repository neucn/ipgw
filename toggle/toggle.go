package toggle

import (
	"ipgw/base"
)

var CmdToggle = &base.Command{
	UsageLine: "ipgw toggle [id | no.] [-v view all]",
	Short:     "切换当前登陆账号",
	Long: `提供多账号切换的功能
使用本功能前需要通过 ipgw login -u xxxxxxxx -p xxxxxx -s 保存过账号
  -v    输出所有中间信息

  ipgw toggle 2018XXXX
    切换到指定学号的保存的账号登陆
  ipgw toggle 2
    切换到第2个保存的账号登陆
  ipgw toggle 2018XXXX -v
    切换到指定学号的保存的账号登陆并输出详细的中间信息
`,
}
