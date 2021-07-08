package message

import (
	"context"
	"fmt"

	"github.com/elap5e/go-mobileqq-api/log"
	"github.com/elap5e/go-mobileqq-api/pb"
)

func SyncUinPairMessages(
	ctx context.Context,
	username string,
	resp *pb.MessageGetMessageResponse,
) {
	for _, uinPairMessage := range resp.GetUinPairMessages() {
		log.Info().
			Uint32("@peer", uinPairMessage.GetPeerUin()).
			Uint32("readAt", uinPairMessage.GetLastReadTime()).
			Msgf("<-> [sync] %d message(s)", len(uinPairMessage.GetMessages()))

		for _, msg := range uinPairMessage.GetMessages() {
			chatID := msg.GetMessageHead().GetGroupInfo().GetGroupCode()
			peerID := uint64(uinPairMessage.GetPeerUin())
			fromID := msg.GetMessageHead().GetFromUin()
			seq := msg.GetMessageHead().GetMessageSeq()

			if msg.GetMessageHead().GetC2CCmd() == 0 {
				chatID = peerID
				peerID = 0
				fromID = 0
			}

			log.Debug().
				Str("@chat", fmt.Sprintf("%d:%d", chatID, peerID)).
				Uint32("@seq", seq).
				Uint64("from", fromID).
				Uint32("time", msg.GetMessageHead().GetMessageTime()).
				Uint32("type", msg.GetMessageHead().GetMessageType()).
				Uint64("uid", msg.GetMessageHead().GetMessageUid()).
				Msg("--> [recv]")
		}
	}
}

func SyncMessage(
	ctx context.Context,
	userID, peerID uint64,
	msg *pb.Message,
) {
	chatID := msg.GetMessageHead().GetGroupInfo().GetGroupCode()
	fromID := msg.GetMessageHead().GetFromUin()
	seq := msg.GetMessageHead().GetMessageSeq()

	log.Debug().
		Str("@chat", fmt.Sprintf("%d:%d", chatID, peerID)).
		Uint32("@seq", seq).
		Uint64("from", fromID).
		Uint32("time", msg.GetMessageHead().GetMessageTime()).
		Uint32("type", msg.GetMessageHead().GetMessageType()).
		Uint64("uid", msg.GetMessageHead().GetMessageUid()).
		Msg("--> [recv]")
}
