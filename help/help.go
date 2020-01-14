package help

import (
	"bufio"
	"io"
	. "ipgw/base"
	"os"
	"strings"
	"text/template"
	"unicode"
	"unicode/utf8"
)

// Help implements the 'help' command.
func Help(w io.Writer, args []string) {
	cmd := Main
Args:
	for i, arg := range args {
		for _, sub := range cmd.Commands {
			if sub.Name == arg {
				cmd = sub
				continue Args
			}
		}

		// helpSuccess is the help command using as many args as possible that would succeed.
		helpSuccess := "ipgw help"
		if i > 0 {
			helpSuccess += " " + strings.Join(args[:i], " ")
		}
		FatalF(CmdNotFound, strings.Join(args, " "), helpSuccess)
	}

	if len(cmd.Commands) > 0 {
		printUsage(w, cmd)
	} else {
		tmpl(w, SimpleUsageTemplate, cmd)
	}
	// not exit 2: succeeded at 'ipgw help cmd'.
	return
}

// An errWriter wraps a writer, recording whether a write error occurred.
type errWriter struct {
	w   io.Writer
	err error
}

func (w *errWriter) Write(b []byte) (int, error) {
	n, err := w.w.Write(b)
	if err != nil {
		w.err = err
	}
	return n, err
}

// tmpl executes the given template text on data, writing the result to w.
func tmpl(w io.Writer, text string, data interface{}) {
	t := template.New("top")
	t.Funcs(template.FuncMap{"trim": strings.TrimSpace, "capitalize": capitalize})
	template.Must(t.Parse(text))
	ew := &errWriter{w: w}
	err := t.Execute(ew, data)
	if ew.err != nil {
		// I/O error writing. Ignore write on closed pipe.
		if strings.Contains(ew.err.Error(), "pipe") {
			os.Exit(1)
		}
		FatalF("writing output: %v", ew.err)
	}
	if err != nil {
		panic(err)
	}
}

func capitalize(s string) string {
	if s == "" {
		return s
	}
	r, n := utf8.DecodeRuneInString(s)
	return string(unicode.ToTitle(r)) + s[n:]
}

func PrintNotFound(cmdName string) {
	FatalF(CmdNotFound, cmdName, "ipgw help")
}

// 在别的包被调用肯定是因为报错，所以公开的方法直接指定为os.Stderr
// 且直接终止程序。
// 用于打印有子命令的命令
func PrintUsage(cmd *Command) {
	printUsage(os.Stderr, cmd)
	os.Exit(2)
}

// 用于打印无子命令的命令
func PrintSimpleUsage(cmd *Command) {
	printSimpleUsage(os.Stderr, cmd)
	os.Exit(2)
}

func printUsage(w io.Writer, cmd *Command) {
	bw := bufio.NewWriter(w)
	tmpl(bw, UsageTemplate, cmd)
	_ = bw.Flush()
}

func printSimpleUsage(w io.Writer, cmd *Command) {
	bw := bufio.NewWriter(w)
	tmpl(bw, SimpleUsageTemplate, cmd)
	_ = bw.Flush()
}
