// 存放命令的结构体与相关函数

package base

import (
	"flag"
	"os"
)

type Command struct {
	// Run runs the command.
	// The args are the arguments after the command name.
	Run func(cmd *Command, args []string)

	Name string

	// UsageLine is the one-line usage message.
	// The words between "ipgw" and the first flag or argument in the line are taken to be the command name.
	UsageLine string

	// Short is the short description shown in the 'ipgw help' output.
	Short string

	// Long is the long message shown in the 'ipgw help <this-command>' output.
	Long string

	// Flag is a set of flags specific to this command.
	Flag flag.FlagSet

	// CustomFlags indicates that the command will do its own
	// flag parsing.
	CustomFlags bool

	// Commands lists the available commands and help topics.
	// The order here is the order in which they are printed by 'go help'.
	// Note that subcommands are in general best avoided.
	Commands []*Command
}

func (c *Command) Usage() {
	ErrorF(CmdUsage, c.UsageLine)
	ErrorF(CmdSeeDetail, c.Name)
	os.Exit(2)
}

// Runnable reports whether the command can be run; otherwise
// it is a documentation pseudo-command such as importpath.
func (c *Command) Runnable() bool {
	return c.Run != nil
}

// Initialize Main Command
var Main = &Command{
	UsageLine: "ipgw",
	Long:      Title,
}
