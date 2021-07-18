package log

import (
	"fmt"
	"time"
)

func PrintMessage(sec int64, peerName, userName, fromName string, peerID, userID, fromID int64, seq int32, text string) {
	if peerID == 0 {
		peerName = userName
	}
	fmt.Println(
		Colorize(time.Unix(sec, 0).Format("[15:04:05]"), ColorBrightBlack, false) +
			Colorize(fmt.Sprintf("[%s(@%du%d:%d)]", peerName, peerID, userID, seq), ColorWhite, false) + " " +
			Colorize(Colorize(fmt.Sprintf("%s(%d):", fromName, fromID), ColorBrightCyan, false), ColorBold, false) + " " +
			Colorize(Colorize(text, ColorBrightWhite, false), ColorBold, false),
	)
}

func PrintMessageSimple(sec int64, peerName, fromName string, seq int32, text string) {
	fmt.Println(
		Colorize(time.Unix(sec, 0).Format("[3:04PM]"), ColorBrightBlack, false) +
			Colorize(fmt.Sprintf("[%s(%d)]", peerName, seq), ColorWhite, false) + " " +
			Colorize(Colorize(fmt.Sprintf("%s:", fromName), ColorBrightCyan, false), ColorBold, false) + " " +
			Colorize(Colorize(text, ColorBrightWhite, false), ColorBold, false),
	)
}
