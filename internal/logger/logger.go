package logger

import (
	"fmt"
)

const (
	ColorReset  = "\033[0m"
	ColorRed    = "\033[31m"
	ColorGreen  = "\033[32m"
	ColorYellow = "\033[33m"
	ColorBlue   = "\033[34m"
	ColorCyan   = "\033[36m"
)

type LogLabel int

const (
	INFO LogLabel = iota
	PASSED
	FAILED
	WARN
	ERROR
	DEBUG
)

func log(level LogLabel, format string, a ...any) {
	msg := fmt.Sprintf(format, a...)

	var color string
	var levelStr string

	switch level {
	case INFO:
		color = ColorGreen
		levelStr = "INFO"
	case PASSED:
		color = ColorGreen
		levelStr = "PASSED"
	case WARN:
		color = ColorYellow
		levelStr = "WARN"
	case FAILED:
		color = ColorRed
		levelStr = "FAILED"
	case ERROR:
		color = ColorRed
		levelStr = "ERROR"
	case DEBUG:
		color = ColorCyan
		levelStr = "DEBUG"
	}

	fmt.Printf("%s[%s] %s%s\n", color, levelStr, ColorReset, msg)
}

func Info(format string, a ...any) {
	log(INFO, format, a...)
}

func Warn(format string, a ...any) {
	log(WARN, format, a...)
}

func Error(format string, a ...any) {
	log(ERROR, format, a...)
}

func Debug(format string, a ...any) {
	log(DEBUG, format, a...)
}

func Passed(format string, a ...any) {
	log(PASSED, format, a...)
}

func Failed(format string, a ...any) {
	log(FAILED, format, a...)
}
