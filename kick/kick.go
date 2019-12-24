package kick

import (
	"fmt"
	"ipgw/base"
	"ipgw/base/cfg"
	"os"
	"regexp"
)

var CmdKick = &base.Command{
	UsageLine: "ipgw kick [-v full view] sid1 sid2 sid3 ...",
	Short:     "使指定设备下线",
	Long: `提供使任意指定设备下线的功能，可同时指定多个
  -v    输出所有中间信息

  ipgw kick XXXXXXX YYYYYYYY
    使指定SID的设备下线
  ipgw kick -v XXXXXXX 
    使指定SID的设备下线并输出详细的中间信息
`,
}

func init() {
	CmdKick.Flag.BoolVar(&cfg.FullView, "v", false, "")

	CmdKick.Run = runKick
}

func runKick(cmd *base.Command, args []string) {
	if len(args) == 0 {
		cmd.Usage()
		return
	}

	pattern := "^\\d{8}$"
	for _, sid := range args {
		matched, _ := regexp.MatchString(pattern, sid)
		if len(sid) != 8 || !matched {
			fmt.Fprintf(os.Stderr, wrongSIDFormat, sid)
			continue
		}
		kickWithSID(sid)
	}
}
