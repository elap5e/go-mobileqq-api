package log

import (
	"fmt"
	"time"
)

func PrintMessage(time time.Time, chatName, peerName, fromName string, chatID, peerID, fromID uint64, seq uint32, text string) {
	if chatID == 0 {
		chatName = peerName
	}
	fmt.Println(
		Colorize(time.Format("[15:04:05]"), ColorBrightBlack, false) +
			Colorize(fmt.Sprintf("[%s(%d:%d:%d)]", chatName, chatID, peerID, seq), ColorWhite, false) + " " +
			Colorize(Colorize(fmt.Sprintf("%s(%d):", fromName, fromID), ColorBrightCyan, false), ColorBold, false) + " " +
			Colorize(Colorize(text, ColorBrightWhite, false), ColorBold, false),
	)
}

func PrintMessageSimple(time time.Time, chatName, peerName, fromName string, chatID, peerID, fromID uint64, seq uint32, text string) {
	if chatID == 0 {
		chatName = peerName
	}
	fmt.Println(
		Colorize(time.Format("[15:04:05]"), ColorBrightBlack, false) +
			Colorize(fmt.Sprintf("[%s(%d:%d)]", chatName, chatID, peerID), ColorWhite, false) + " " +
			Colorize(Colorize(fmt.Sprintf("%s:", fromName), ColorBrightCyan, false), ColorBold, false) + " " +
			Colorize(Colorize(text, ColorBrightWhite, false), ColorBold, false),
	)
}
