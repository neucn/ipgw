package api

import (
	"ipgw/api/code"
	"ipgw/api/login"
	"ipgw/api/proxy"
	. "ipgw/base"
	"ipgw/lib"
	"os"
)

var CmdAPI = &Command{
	UsageLine: "ipgw api",
	Short:     "API",
	Long:      `To Be Done.`,
}

var (
// 错误码 = ${出错命名空间}*10 + ${具体错误码}
// 若具体错误码>10，则占据前一个或后一个命名空间
// 若无错误则输出 0
// 除0以外，其他具体错误码都不能是10的倍数，例如login的具体错误码应该从1开始排
/*
	命令码:
		0 - 用于基本具体错误码
		1 - login about
		2 - proxy about


	基本具体错误码:
		0 - 成功无错误
		1 - 无该命令
		2 - 参数缺失
		3 - 版本不兼容（例如本工具进行了接口更新，无法兼容基于老版本接口的tool，则返回错误使tool感知
		4 - 配置文件读取失败
		5 - 网络连接错误
		更多等待补充

	个别具体错误码:
		详见各api实现中的注释
*/

)

func init() {
	CmdAPI.Run = runAPI
}

func runAPI(cmd *Command, args []string) {
	// args [要求版本 命令 参数一 参数二]
	if len(args) < 2 {
		// 无参数直接结束程序
		Fatal(code.GlobalMissArgs)
	}

	// 对args[0]进行版本判断
	if !lib.IsAPICompatible(args[0]) {
		Fatal(code.GlobalNotCompatible)
	}

	// 命令转发
	args = args[1:]
	switch args[0] {
	case "login":
		login.Login.Run(login.Login, args[1:])
	case "proxy":
		proxy.Proxy.Run(proxy.Proxy, args[1:])
	default:
		Fatal(code.GlobalNoCmd)
	}

	// 进行到这说明执行成功，因此输出错误码0
	Error(code.GlobalSuccess)
	os.Exit(0)
}
