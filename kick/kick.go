package kick

import (
	. "ipgw/base"
	"ipgw/ctx"
	"regexp"
)

var CmdKick = &Command{
	Name:      "kick",
	UsageLine: "ipgw kick [-v view all] sid1 sid2 sid3 ...",
	Short:     "使指定设备下线",
	Long: `提供使任意指定设备下线的功能，可同时指定多个
  -v    输出所有中间信息

  ipgw kick xxxxxxxx yyyyyyyy
    使指定SID的设备下线
  ipgw kick -v xxxxxxxx
    使指定SID的设备下线并输出详细的中间信息
`,
}

func init() {
	CmdKick.Flag.BoolVar(&ctx.FullView, "v", false, "")

	CmdKick.Run = runKick
}

func runKick(cmd *Command, args []string) {
	// 若无参则打印使用说明并结束
	if len(args) == 0 {
		cmd.Usage()
		return
	}

	// 正则匹配参数是否符合SID格式
	pattern := "^\\d{8}$"
	c := ctx.NewCtx()

	for _, sid := range args {
		matched, _ := regexp.MatchString(pattern, sid)
		if !matched {
			// 格式错误，跳到下一个SID
			ErrorF(wrongSIDFormat, sid)
			continue
		}
		// Kick，因为封装所以逻辑有问题，见函数里的说明
		kick(c, sid)
	}
}
