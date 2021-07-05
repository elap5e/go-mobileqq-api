package log

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const (
	ColorBold = 1
)

const (
	ColorBlack = iota + 30
	ColorRed
	ColorGreen
	ColorYellow
	ColorBlue
	ColorMagenta
	ColorCyan
	ColorWhite
)

const (
	ColorBrightBlack = iota + 90
	ColorBrightRed
	ColorBrightGreen
	ColorBrightYellow
	ColorBrightBlue
	ColorBrightMagenta
	ColorBrightCyan
	ColorBrightWhite
)

// var consoleTimeFormat = time.RFC3339Nano

var consoleTimeFormat = "2006-01-02T15:04:05.000Z07:00"

func Colorize(s interface{}, c int) string {
	return fmt.Sprintf("\x1b[%dm%v\x1b[0m", c, s)
}

var consoleFormatLevel = func(i interface{}) string {
	var l string
	if ll, ok := i.(string); ok {
		switch ll {
		case "trace":
			l = Colorize("[TRAC]", ColorBrightWhite)
		case "debug":
			l = Colorize("[DEBU]", ColorBrightWhite)
		case "info":
			l = Colorize("[INFO]", ColorBrightCyan)
		case "warn":
			l = Colorize("[WARN]", ColorBrightYellow)
		case "error":
			l = Colorize(Colorize("[ERRO]", ColorBrightRed), ColorBold)
		case "fatal":
			l = Colorize(Colorize("[FATA]", ColorBrightRed), ColorBold)
		case "panic":
			l = Colorize(Colorize("[PANI]", ColorBrightRed), ColorBold)
		default:
			l = Colorize("[????]", ColorBold)
		}
	} else {
		if i == nil {
			l = Colorize("[????]", ColorBold)
		} else {
			l = strings.ToUpper(fmt.Sprintf("%s", i))[0:3]
		}
	}
	return l
}

var consoleFormatCaller = func(i interface{}) string {
	var c string
	if cc, ok := i.(string); ok {
		c = cc
	}
	if len(c) > 0 {
		if cwd, err := os.Getwd(); err == nil {
			if rel, err := filepath.Rel(cwd, c); err == nil {
				c = rel
			}
		}
		c = Colorize(c, ColorBold) + Colorize(" >", ColorCyan)
	}
	return c
}

var consoleFormatMessage = func(i interface{}) string {
	if i == nil {
		return ""
	}
	return fmt.Sprintf("%s", i)
}

var consoleFormatFieldName = func(i interface{}) string {
	return Colorize(fmt.Sprintf("%s:", i), ColorBrightCyan)
}

var consoleFormatFieldValue = func(i interface{}) string {
	return fmt.Sprintf("%s", i)
}

var consoleFormatErrFieldName = func(i interface{}) string {
	return Colorize(Colorize(fmt.Sprintf("%s:", i), ColorBrightRed), ColorBold)
}

var consoleFormatErrFieldValue = func(i interface{}) string {
	return Colorize(Colorize(fmt.Sprintf("%s", i), ColorBrightRed), ColorBold)
}
