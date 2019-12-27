package main

import (
	"flag"
	"fmt"
	"ipgw/base"
	"ipgw/base/cfg"
	"ipgw/fix"
	"ipgw/help"
	"ipgw/kick"
	"ipgw/list"
	"ipgw/login"
	"ipgw/logout"
	"ipgw/test"
	"ipgw/text"
	"ipgw/update"
	"ipgw/version"
	"log"
	"os"
	"strings"
)

func init() {
	base.IPGW.Commands = []*base.Command{
		version.CmdVersion,
		login.CmdLogin,
		logout.CmdLogout,
		//toggle.CmdToggle,
		list.CmdList,
		kick.CmdKick,
		test.CmdTest,
		fix.CmdFix,
		update.CmdUpdate,
	}
	base.Usage = mainUsage
}

func main() {
	flag.Usage = base.Usage
	flag.Parse()
	log.SetFlags(0)

	args := flag.Args()
	if len(args) < 1 {
		login.CmdLogin.Run(login.CmdLogin, nil)
		os.Exit(2)
	}

	cfg.CmdName = args[0] // for error messages

	// 处理help
	if args[0] == "help" {
		help.Help(os.Stdout, args[1:])
		return
	}

BigCmdLoop:
	for bigCmd := base.IPGW; ; {
		for _, cmd := range bigCmd.Commands {
			if cmd.Name() != args[0] {
				continue
			}
			if len(cmd.Commands) > 0 {
				bigCmd = cmd
				args = args[1:]
				if len(args) == 0 {
					help.PrintUsage(os.Stderr, bigCmd)
					base.SetExitStatus(2)
					base.Exit()
				}
				if args[0] == "help" {
					help.Help(os.Stdout, append(strings.Split(cfg.CmdName, " "), args[1:]...))
					return
				}
				cfg.CmdName += " " + args[0]
				continue BigCmdLoop
			}
			if !cmd.Runnable() {
				continue
			}
			cmd.Flag.Usage = func() { cmd.Usage() }
			if cmd.CustomFlags {
				args = args[1:]
			} else {
				cmd.Flag.Parse(args[1:])
				args = cmd.Flag.Args()
			}
			cmd.Run(cmd, args)
			base.Exit()
			return
		}
		fmt.Fprintf(os.Stderr, text.HelpNotFound, cfg.CmdName, "ipgw help")
		base.SetExitStatus(2)
		base.Exit()
	}
}

func mainUsage() {
	help.PrintUsage(os.Stderr, base.IPGW)
	os.Exit(2)
}
