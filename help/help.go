package help

import (
	"bufio"
	"fmt"
	"io"
	"ipgw/base"
	"ipgw/text"
	"os"
	"strings"
	"text/template"
	"unicode"
	"unicode/utf8"
)

// Help implements the 'help' command.
func Help(w io.Writer, args []string) {
	cmd := base.IPGW
Args:
	for i, arg := range args {
		for _, sub := range cmd.Commands {
			if sub.Name() == arg {
				cmd = sub
				continue Args
			}
		}

		// helpSuccess is the help command using as many args as possible that would succeed.
		helpSuccess := "ipgw help"
		if i > 0 {
			helpSuccess += " " + strings.Join(args[:i], " ")
		}
		fmt.Fprintf(os.Stderr, text.HelpNotFound, strings.Join(args, " "), helpSuccess)
		base.SetExitStatus(2) // failed at 'ipgw help cmd'
		base.Exit()
	}

	if len(cmd.Commands) > 0 {
		PrintUsage(os.Stdout, cmd)
	} else {
		tmpl(os.Stdout, text.HelpTemplate, cmd)
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
			base.SetExitStatus(1)
			base.Exit()
		}
		base.Fatalf("writing output: %v", ew.err)
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

func PrintUsage(w io.Writer, cmd *base.Command) {
	bw := bufio.NewWriter(w)
	tmpl(bw, text.HelpUsageTemplate, cmd)
	bw.Flush()
}
