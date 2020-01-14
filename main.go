package main

import (
	"flag"
	"ipgw/api"
	. "ipgw/base"
	"ipgw/fix"
	"ipgw/help"
	"ipgw/kick"
	"ipgw/list"
	"ipgw/login"
	"ipgw/logout"
	"ipgw/test"
	"ipgw/tool"
	"ipgw/update"
	"ipgw/version"
	"os"
	"strings"
)

func init() {
	Main.Commands = []*Command{
		login.CmdLogin,
		logout.CmdLogout,
		kick.CmdKick,
		list.CmdList,
		tool.CmdTool,
		test.CmdTest,
		update.CmdUpdate,
		fix.CmdFix,
		version.CmdVersion,
	}
}

func main() {
	flag.Usage = func() { help.PrintUsage(Main) }

	// 第一次解析
	flag.Parse()

	// 获取命令行参数列表
	args := flag.Args()

	// 实现`ipgw`直接登陆
	if len(args) < 1 {
		login.CmdLogin.Run(login.CmdLogin, nil)
		return
	}

	// 处理help命令
	if args[0] == "help" {
		help.Help(os.Stdout, args[1:])
		return
	}

	// 处理api命令
	if args[0] == "api" {
		api.CmdAPI.Run(api.CmdAPI, args[1:])
		return
	}

	// 解析
	parse(args)
}

// 主体循环解析
func parse(args []string) {
	cmdName := args[0] // for error messages
BigCmdLoop:
	for bigCmd := Main; ; {
		for _, cmd := range bigCmd.Commands {
			if cmd.Name != args[0] {
				continue
			}
			if len(cmd.Commands) > 0 {
				bigCmd = cmd
				args = args[1:]
				if len(args) == 0 {
					help.PrintUsage(bigCmd)
				}
				if args[0] == "help" {
					help.Help(os.Stdout, append(strings.Split(cmdName, " "), args[1:]...))
					return
				}
				cmdName += " " + args[0]
				continue BigCmdLoop
			}
			if !cmd.Runnable() {
				continue
			}
			cmd.Flag.Usage = func() { cmd.Usage() }
			if cmd.CustomFlags {
				args = args[1:]
			} else {
				_ = cmd.Flag.Parse(args[1:])
				args = cmd.Flag.Args()
			}
			cmd.Run(cmd, args)
			// 假设能到达这里的业务都是正常工作的，因此以0退出
			os.Exit(0)
			return
		}
		help.PrintNotFound(cmdName)
	}
}
