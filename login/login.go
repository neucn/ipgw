package login

import (
	"fmt"
	"ipgw/base"
	"ipgw/base/cfg"
)

var (
	u, p, c string

	s bool
)

func init() {
	CmdLogin.Flag.StringVar(&u, "u", "", "")
	CmdLogin.Flag.StringVar(&p, "p", "", "")
	CmdLogin.Flag.StringVar(&c, "c", "", "")

	CmdLogin.Flag.BoolVar(&s, "s", false, "")
	CmdLogin.Flag.BoolVar(&cfg.FullView, "v", false, "")

	CmdLogin.Run = runLogin
}

var CmdLogin = &base.Command{
	UsageLine: "ipgw login [-u username] [-p password] [-s save] [-c cookie] [-v full view] ",
	Short:     "基础登陆",
	Long: `提供登陆校园网关功能
  -u    登陆账户
  -p    登陆密码
  -s    保存该账户
  -c    使用cookie登陆
  -v    输出所有中间信息

  ipgw login -u 学号 -p 密码
    使用指定账号登陆网关
  ipgw login -u 学号 -p 密码 -s
    若在登陆时开启-s, 本次登陆的账号信息将被保存在用户目录下的.ipgw文件中
  ipgw login
    在已经使用-s保存了账户信息的情况下，可以直接使用已经保存的账号登录
  ipgw
    [推荐] 是的没错，在已经使用-s保存了账号信息的情况下，直接执行ipgw即可完成网关登陆
  ipgw login -c "ST-XXXXXX-XXXXXXXXXXXXXXXXXXXX-tpass"
    使用指定cookie登陆
  ipgw login [arguments] -v
    打印登陆过程中的每一步信息
`,
}

func init() {
	CmdLogin.Run = runLogin // break init cycle
}

func runLogin(cmd *base.Command, args []string) {
	fmt.Println(u, p, s, c)
}
