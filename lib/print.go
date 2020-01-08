// 统一的输出

package lib

import (
	"fmt"
	"os"
)

func Error(msg string) {
	_, _ = fmt.Fprintln(os.Stderr, msg)
}

func ErrorF(format string, a ...interface{}) {
	_, _ = fmt.Fprintf(os.Stderr, format, a...)
}

func Fatal(msg string) {
	Error(msg)
	os.Exit(2)
}

func FatalF(format string, a ...interface{}) {
	ErrorF(format, a...)
	os.Exit(2)
}

func Info(a ...interface{}) {
	fmt.Print(a...)
}

func InfoLine(a ...interface{}) {
	fmt.Println(a...)
}

func InfoF(format string, a ...interface{}) {
	fmt.Printf(format, a...)
}
