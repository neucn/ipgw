package api

import (
	"ipgw/base"
	. "ipgw/lib"
)

var API = &base.Command{
	UsageLine: "ipgw api",
	Short:     "API",
	Long:      `To Be Done.`,
}

var (
// 错误输出格式: `出错命令码 具体错误码`
/*
	出错命令码:
		0 - 用于通用具体错误码
		1 - api login

	具体错误码分类:
		负数代表内部错误	如: 网络请求错误，程序运行错误
		0代表参数缺失
		正数代表业务错误	如：账号密码错误，Cookie失效

	通用具体错误码:
		 0 - 参数缺失
		-1 - 版本不兼容（例如本工具进行了接口更新，无法兼容基于老版本接口的tool，则返回错误使tool感知
		-2 - 配置文件读取失败
		-3 - 网络连接错误
		更多等待补充

	个别具体错误码:
		详见各api实现中的注释
*/

)

func init() {
	API.Run = runAPI
}

func runAPI(cmd *base.Command, args []string) {
	// args [要求版本 命令 参数一 参数二]
	if len(args) < 2 {
		// 无参数直接结束程序
		Fatal(globalMissArgs)
	}

	// 对args[0]进行版本判断
	if !isCompatible(args[0]) {
		Fatal(globalNotCompatible)
	}

	// 命令转发，目前只有login需求
	args = args[1:]
	switch args[0] {
	case "login":
		Login.Run(Login, args[1:])
	}

}

// 检验是否兼容
// 暂时简单判断
func isCompatible(v string) bool {
	return base.API == v
}
