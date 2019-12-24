package text

var (
	HelpNotFound  = "无此命令: ipgw %s\n请使用 %s 查看帮助信息\n"
	HelpSeeDetail = "使用 'ipgw help %s' 获取命令的完整帮助信息。\n"
	HelpUsage     = "用法: %s\n"

	HelpUsageTemplate = `{{.Long | trim}}

用法:

	{{.UsageLine}} <command> [arguments]

命令:
{{range .Commands}}{{if or (.Runnable) .Commands}}
	{{.Name | printf "%-11s"}} {{.Short}}{{end}}{{end}}

使用 "ipgw help{{with .LongName}} {{.}}{{end}} <command>" 获取某个命令的完整帮助信息。
{{if eq (.UsageLine) "ipgw"}}
有意见或建议欢迎发送至邮箱 i@shangyes.net
{{end}}
`

	HelpTemplate = `{{if .Runnable}}用法: {{.UsageLine}}

{{end}}{{.Long | trim}}
`
)
