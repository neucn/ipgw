package test

import "ipgw/base"

var CmdTestSpeed = &base.Command{
	UsageLine: "ipgw test speed [-v full view]",
	Short:     "校园网测速",
	Long: `提供对于校园网的测速功能
  -v    输出所有中间信息

  ipgw test speed
    校园网测速
  ipgw test speed -v
    校园网测速并输出详细中间信息
`,
}
