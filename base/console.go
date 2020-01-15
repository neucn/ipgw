// 统一的输出

package base

import (
	"fmt"
	"os"
)

func Error(msg ...interface{}) {
	_, _ = fmt.Fprint(os.Stderr, msg...)
}

func ErrorL(msg ...interface{}) {
	_, _ = fmt.Fprintln(os.Stderr, msg...)
}

func ErrorF(format string, msg ...interface{}) {
	_, _ = fmt.Fprintf(os.Stderr, format, msg...)
}

func Fatal(msg ...interface{}) {
	_, _ = fmt.Fprint(os.Stderr, msg...)
	os.Exit(2)
}

func FatalL(msg ...interface{}) {
	_, _ = fmt.Fprintln(os.Stderr, msg...)
	os.Exit(2)
}

func FatalF(format string, msg ...interface{}) {
	_, _ = fmt.Fprintf(os.Stderr, format, msg...)
	os.Exit(2)
}

func Info(msg ...interface{}) {
	fmt.Print(msg...)
}

func InfoL(msg ...interface{}) {
	fmt.Println(msg...)
}

func InfoF(format string, msg ...interface{}) {
	fmt.Printf(format, msg...)
}
