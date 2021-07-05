package log

import (
	"fmt"
)

const (
	msgDumpRecvInfo  = ">>> [dump] "
	msgDumpSendInfo  = "<<< [dump] "
	msgInitInfo      = "··· [init] "
	msgInitError     = "x_x [init] "
	msgConnRecvWarn  = "->· [conn] "
	msgConnRecvError = "->x [conn] "
)

var MsgConnRecvWarn = func(i interface{}) string {
	return Colorize(fmt.Sprintf(msgConnRecvWarn+"%s", i), ColorBrightYellow)
}

var MsgConnRecvWarnf = func(format string, v ...interface{}) string {
	return Colorize(fmt.Sprintf(msgConnRecvWarn+format, v...), ColorBrightYellow)
}

var MsgConnRecvError = func(i interface{}) string {
	return Colorize(fmt.Sprintf(msgConnRecvError+"%s", i), ColorBrightYellow)
}

var MsgConnRecvErrorf = func(format string, v ...interface{}) string {
	return Colorize(fmt.Sprintf(msgConnRecvError+format, v...), ColorBrightYellow)
}

var MsgDumpRecvInfo = func(i interface{}) string {
	return Colorize(fmt.Sprintf(msgDumpRecvInfo+"%s", i), ColorBrightWhite)
}

var MsgDumpRecvInfof = func(format string, v ...interface{}) string {
	return Colorize(fmt.Sprintf(msgDumpRecvInfo+format, v...), ColorBrightWhite)
}

var MsgDumpSendInfo = func(i interface{}) string {
	return Colorize(fmt.Sprintf(msgDumpSendInfo+"%s", i), ColorBrightWhite)
}

var MsgDumpSendInfof = func(format string, v ...interface{}) string {
	return Colorize(fmt.Sprintf(msgDumpSendInfo+format, v...), ColorBrightWhite)
}

var MsgInitInfo = func(i interface{}) string {
	return fmt.Sprintf(msgInitInfo+"%s", i)
}

var MsgInitInfof = func(format string, v ...interface{}) string {
	return fmt.Sprintf(msgInitInfo+format, v...)
}

var MsgInitError = func(i interface{}) string {
	return Colorize(Colorize(fmt.Sprintf(msgInitError+"%s", i), ColorBrightRed), ColorBold)
}

var MsgInitErrorf = func(format string, v ...interface{}) string {
	return Colorize(Colorize(fmt.Sprintf(msgInitError+format, v...), ColorBrightRed), ColorBold)
}
