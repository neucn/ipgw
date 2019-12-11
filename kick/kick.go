package kick

import "C"
import (
	"fmt"
	"ipgw/base"
	"ipgw/base/cfg"
)

var CmdKick = &base.Command{
	UsageLine: "ipgw kick [-d device] [-s ssid] [-v full view]",
	Short:     "使指定设备下线",
	Long: `提供使任意指定设备下线的功能
  -d    指定设备
  -s	指定ssid
  -v    输出所有中间信息

  ipgw kick -s XXXXXXX
    使指定SSID的设备下线
  ipgw kick -d 2
    使当前账号的第二个设备下线，当前账号的设备可由 ipgw list -d 查看
  ipgw kick -s XXXXXXX -v
    使指定SSID的设备下线并输出详细的中间信息
`,
}

var (
	d, s string
)

func init() {
	CmdKick.Flag.StringVar(&d, "d", "", "")
	CmdKick.Flag.StringVar(&s, "s", "", "")

	CmdKick.Flag.BoolVar(&cfg.FullView, "v", false, "")

	CmdKick.Run = runKick
}

func runKick(cmd *base.Command, args []string) {
	fmt.Println(d, s)
}
