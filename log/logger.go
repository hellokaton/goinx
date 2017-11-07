package log

import (
	"fmt"
	"time"
)

var LogLevelNum = 2

const (
	color_red = uint8(iota + 91)
	color_green
	color_yellow
	color_blue
	color_magenta //洋红

	info  = "[INFO]"
	debug = "[DEBUG]"
	erro  = "[ERROR]"
	warn  = "[WARN]"
)

// see complete color rules in document in https://en.wikipedia.org/wiki/ANSI_escape_code#cite_note-ecma48-13
func Debug(format string, a ...interface{}) {
	if LogLevelNum <= 1 {
		prefix := yellow(debug)
		fmt.Println(formatLog(prefix), fmt.Sprintf(format, a...))
	}
}

func Info(format string, a ...interface{}) {
	if LogLevelNum <= 2 {
		prefix := green(info)
		fmt.Println(formatLog(prefix), fmt.Sprintf(format, a...))
	}
}

// func Success(format string, a ...interface{}) {
// 	prefix := blue(info)
// 	fmt.Println(formatLog(prefix), fmt.Sprintf(format, a...))
// }

func Warning(format string, a ...interface{}) {
	if LogLevelNum <= 3 {
		prefix := magenta(warn)
		fmt.Println(formatLog(prefix), fmt.Sprintf(format, a...))
	}
}

func Error(format string, a ...interface{}) {
	if LogLevelNum <= 4 {
		prefix := red(erro)
		fmt.Println(formatLog(prefix), fmt.Sprintf(format, a...))
	}
}

func red(s string) string {
	return fmt.Sprintf("\x1b[%dm%s\x1b[0m", color_red, s)
}

func green(s string) string {
	return fmt.Sprintf("\x1b[%dm%s\x1b[0m", color_green, s)
}

func yellow(s string) string {
	return fmt.Sprintf("\x1b[%dm%s\x1b[0m", color_yellow, s)
}

func blue(s string) string {
	return fmt.Sprintf("\x1b[%dm%s\x1b[0m", color_blue, s)
}

func magenta(s string) string {
	return fmt.Sprintf("\x1b[%dm%s\x1b[0m", color_magenta, s)
}

func formatLog(prefix string) string {
	return time.Now().Format("2006/01/02 15:04:05") + " " + prefix
}
