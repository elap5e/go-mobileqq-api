package log

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/rs/zerolog"
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

var consoleTimeFormat = "2006-01-02T15:04:05.000Z07:00"

var ConsoleWriter = zerolog.ConsoleWriter{
	Out:                 os.Stdout,
	NoColor:             false,
	TimeFormat:          consoleTimeFormat,
	FormatLevel:         consoleFormatLevel(false),
	FormatCaller:        consoleFormatCaller(false),
	FormatMessage:       consoleFormatMessage,
	FormatFieldName:     consoleFormatFieldName(false),
	FormatFieldValue:    consoleFormatFieldValue,
	FormatErrFieldName:  consoleFormatErrFieldName(false),
	FormatErrFieldValue: consoleFormatErrFieldValue(false),
}

var ConsoleWriterNoColor = zerolog.ConsoleWriter{
	Out:                 os.Stdout,
	NoColor:             true,
	TimeFormat:          consoleTimeFormat,
	FormatLevel:         consoleFormatLevel(true),
	FormatCaller:        consoleFormatCaller(true),
	FormatMessage:       consoleFormatMessage,
	FormatFieldName:     consoleFormatFieldName(true),
	FormatFieldValue:    consoleFormatFieldValue,
	FormatErrFieldName:  consoleFormatErrFieldName(true),
	FormatErrFieldValue: consoleFormatErrFieldValue(true),
}

func Colorize(s interface{}, c int, noColor bool) string {
	if noColor {
		return fmt.Sprintf("%s", s)
	}
	return fmt.Sprintf("\x1b[%dm%v\x1b[0m", c, s)
}

func consoleFormatLevel(noColor bool) zerolog.Formatter {
	return func(i interface{}) string {
		var l string
		if ll, ok := i.(string); ok {
			switch ll {
			case "trace":
				l = Colorize("[TRAC]", ColorBrightWhite, noColor)
			case "debug":
				l = Colorize("[DEBU]", ColorBrightWhite, noColor)
			case "info":
				l = Colorize("[INFO]", ColorBrightCyan, noColor)
			case "warn":
				l = Colorize("[WARN]", ColorBrightYellow, noColor)
			case "error":
				l = Colorize(Colorize("[ERRO]", ColorBrightRed, noColor), ColorBold, noColor)
			case "fatal":
				l = Colorize(Colorize("[FATA]", ColorBrightRed, noColor), ColorBold, noColor)
			case "panic":
				l = Colorize(Colorize("[PANI]", ColorBrightRed, noColor), ColorBold, noColor)
			default:
				l = Colorize("[????]", ColorBold, noColor)
			}
		} else {
			if i == nil {
				l = Colorize("[????]", ColorBold, noColor)
			} else {
				l = strings.ToUpper(fmt.Sprintf("%s", i))[0:3]
			}
		}
		return l
	}
}

func consoleFormatCaller(noColor bool) zerolog.Formatter {
	return func(i interface{}) string {
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
			c = Colorize(c, ColorBold, noColor) + Colorize(" >", ColorCyan, noColor)
		}
		return c
	}
}

var consoleFormatMessage = func(i interface{}) string {
	if i == nil {
		return ""
	}
	return fmt.Sprintf("%s", i)
}

func consoleFormatFieldName(noColor bool) zerolog.Formatter {
	return func(i interface{}) string {
		return Colorize(fmt.Sprintf("%s:", i), ColorBrightCyan, noColor)
	}
}

var consoleFormatFieldValue = func(i interface{}) string {
	return fmt.Sprintf("%s", i)
}

func consoleFormatErrFieldName(noColor bool) zerolog.Formatter {
	return func(i interface{}) string {
		return Colorize(Colorize(fmt.Sprintf("%s:", i), ColorBrightRed, noColor), ColorBold, noColor)
	}
}

func consoleFormatErrFieldValue(noColor bool) zerolog.Formatter {
	return func(i interface{}) string {
		return Colorize(Colorize(fmt.Sprintf("%s", i), ColorBrightRed, noColor), ColorBold, noColor)
	}
}
