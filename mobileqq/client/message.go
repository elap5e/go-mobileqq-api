package client

import (
	"fmt"

	"github.com/elap5e/go-mobileqq-api/log"
	"github.com/elap5e/go-mobileqq-api/pb"
)

func syncUinPairMessage(uinPairMessage *pb.UinPairMessage) {
	log.Info().
		Uint64("@peer", uinPairMessage.GetPeerUin()).
		Int64("readAt", uinPairMessage.GetLastReadTime()).
		Msgf("<-> [sync] %d message(s)", len(uinPairMessage.GetMessages()))

	for _, msg := range uinPairMessage.GetMessages() {
		chatID := msg.GetMessageHead().GetGroupInfo().GetGroupCode()
		peerID := uint64(0)
		if msg.GetMessageHead().GetC2CCmd() != 0 {
			chatID = msg.GetMessageHead().GetC2CTempMessageHead().GetGroupCode()
			peerID = uinPairMessage.GetPeerUin()
		}

		log.Debug().
			Str("@chat", fmt.Sprintf("%d:%d", chatID, peerID)).
			Uint32("@seq", msg.GetMessageHead().GetMessageSeq()).
			Uint64("from", msg.GetMessageHead().GetFromUin()).
			Int64("time", msg.GetMessageHead().GetMessageTime()).
			Uint32("type", msg.GetMessageHead().GetMessageType()).
			Uint64("uid", msg.GetMessageHead().GetMessageUid()).
			Msg("--> [recv]")
	}
}

func syncMessage(msg *pb.Message) {
	chatID := msg.GetMessageHead().GetGroupInfo().GetGroupCode()
	peerID := uint64(0)
	if msg.GetMessageHead().GetC2CCmd() != 0 {
		chatID = msg.GetMessageHead().GetC2CTempMessageHead().GetGroupCode()
		peerID = msg.GetMessageHead().GetToUin()
	}

	log.Debug().
		Str("@chat", fmt.Sprintf("%d:%d", chatID, peerID)).
		Uint32("@seq", msg.GetMessageHead().GetMessageSeq()).
		Uint64("from", msg.GetMessageHead().GetFromUin()).
		Int64("time", msg.GetMessageHead().GetMessageTime()).
		Uint32("type", msg.GetMessageHead().GetMessageType()).
		Uint64("uid", msg.GetMessageHead().GetMessageUid()).
		Msg("--> [recv]")
}
