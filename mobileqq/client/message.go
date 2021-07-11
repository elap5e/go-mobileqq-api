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
		peerID := msg.GetMessageHead().GetC2CTempMessageHead().GetGroupCode()
		userID := uinPairMessage.GetPeerUin()
		if msg.GetMessageHead().GetC2CCmd() == 0 {
			peerID = uinPairMessage.GetPeerUin()
			userID = uint64(0)
		}

		log.Debug().
			Str("@chat", fmt.Sprintf("@%d_%d", peerID, userID)).
			Uint32("@seq", msg.GetMessageHead().GetMessageSeq()).
			Uint64("from", msg.GetMessageHead().GetFromUin()).
			Int64("time", msg.GetMessageHead().GetMessageTime()).
			Uint32("type", msg.GetMessageHead().GetMessageType()).
			Uint64("uid", msg.GetMessageHead().GetMessageUid()).
			Msg("--> [recv]")
	}
}

func syncMessage(msg *pb.Message) {
	peerID := msg.GetMessageHead().GetC2CTempMessageHead().GetGroupCode()
	userID := msg.GetMessageHead().GetToUin()
	if msg.GetMessageHead().GetC2CCmd() == 0 {
		peerID = msg.GetMessageHead().GetGroupInfo().GetGroupCode()
		userID = uint64(0)
	}

	log.Debug().
		Str("@chat", fmt.Sprintf("@%d_%d", peerID, userID)).
		Uint32("@seq", msg.GetMessageHead().GetMessageSeq()).
		Uint64("from", msg.GetMessageHead().GetFromUin()).
		Int64("time", msg.GetMessageHead().GetMessageTime()).
		Uint32("type", msg.GetMessageHead().GetMessageType()).
		Uint64("uid", msg.GetMessageHead().GetMessageUid()).
		Msg("--> [recv]")
}
