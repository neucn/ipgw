package list

import (
	"fmt"
	"ipgw/base"
	"ipgw/base/cfg"
	"ipgw/text"
	"os"
	"strings"
)

var CmdList = &base.Command{
	CustomFlags: true,
	UsageLine:   "ipgw list [-a all] [-i info] [-d device] [-c cookie] [-s saved account] [-h history] page [-v full view]",
	Short:       "获取各类信息",
	Long: `提供校园网信息查询功能
  -a    列出所有信息
  -c    列出当前cookie
  -d    列出登陆设备
  -h    列出校园网使用日志
  -i    列出账户信息
  -s    列出保存的账户
  -v    输出所有中间信息

  ipgw list -a
    效果等同于 ipgw list -disch 1，列出所有信息
  ipgw list -i
    若本次登陆使用的是本工具，则查看当前账户的具体信息
    若本次登陆不是本工具，但已使用-s保存过账户信息，则查看保存账户的具体信息
  ipgw list -d
    若本次登陆使用的是本工具，则查看当前账户的已登录设备
    若本次登陆不是本工具，但已使用-s保存过账户信息，则查看保存账户的已登录设备
  ipgw list -c
    查看当前储存的有效cookie
  ipgw list -s
    查看当前保存的所有账户
  ipgw list -h 1
    若本次登陆使用的是本工具，则查看当前账户的第一页使用日志
    若本次登陆不是本工具，但已使用-s保存过账户信息，则查看保存账户的第一页使用日志
  ipgw list -av
    打印所有信息且输出每一步的详细信息
`,
}

var (
	flags            = []int32{'a', 'c', 'd', 'h', 'i', 's'}
	a, i, d, c, s, h bool
)

func init() {
	CmdList.Flag.BoolVar(&a, "a", false, "列出所有信息")
	CmdList.Flag.BoolVar(&i, "i", false, "列出账户信息")
	CmdList.Flag.BoolVar(&d, "d", false, "列出登陆设备")
	CmdList.Flag.BoolVar(&c, "c", false, "列出当前cookie")
	CmdList.Flag.BoolVar(&s, "s", false, "列出保存的账户")
	CmdList.Flag.BoolVar(&h, "h", false, "列出校园网使用日志")

	CmdList.Flag.BoolVar(&cfg.FullView, "v", false, "输出所有中间信息")

	CmdList.Run = runList // break init cycle
}

func runList(cmd *base.Command, args []string) {
	parse(cmd, args)
	fmt.Println(a, c, d, h, i, s, cfg.FullView, cmd.Flag.Args())
}

func parse(cmd *base.Command, args []string) {
	separated := make([]string, 0, len(args))
	for _, flagChar := range args {
		if len(flagChar) > 2 && strings.HasPrefix(flagChar, "-") {
		charLoop:
			for _, c := range flagChar[1:] {
				for _, f := range flags {
					if c == f {
						separated = append(separated, "-"+string(c))
						continue charLoop
					}
				}
				fmt.Fprintf(os.Stdout, text.ListArgNotFound, string(c))
				cmd.Usage()
			}
			continue
		}
		separated = append(separated, flagChar)
	}
	cmd.Flag.Parse(separated)
}
