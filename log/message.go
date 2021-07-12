package log

import (
	"fmt"
	"time"
)

func PrintMessage(time time.Time, peerName, userName, fromName string, peerID, userID, fromID uint64, seq uint32, text string) {
	if peerID == 0 {
		peerName = userName
	}
	fmt.Println(
		Colorize(time.Format("[15:04:05]"), ColorBrightBlack, false) +
			Colorize(fmt.Sprintf("[%s(@%d_%d:%d)]", peerName, peerID, userID, seq), ColorWhite, false) + " " +
			Colorize(Colorize(fmt.Sprintf("%s(%d):", fromName, fromID), ColorBrightCyan, false), ColorBold, false) + " " +
			Colorize(Colorize(text, ColorBrightWhite, false), ColorBold, false),
	)
}

func PrintMessageSimple(time time.Time, peerName, fromName string, seq uint32, text string) {
	fmt.Println(
		Colorize(time.Format("[3:04PM]"), ColorBrightBlack, false) +
			Colorize(fmt.Sprintf("[%s(%d)]", peerName, seq), ColorWhite, false) + " " +
			Colorize(Colorize(fmt.Sprintf("%s:", fromName), ColorBrightCyan, false), ColorBold, false) + " " +
			Colorize(Colorize(text, ColorBrightWhite, false), ColorBold, false),
	)
}
