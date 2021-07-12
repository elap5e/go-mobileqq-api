package client

import (
	"fmt"
	"strconv"
	"time"

	"github.com/elap5e/go-mobileqq-api/log"
	"github.com/elap5e/go-mobileqq-api/pb"
	"github.com/rs/zerolog"
)

func (c *Client) PrintMessage(time time.Time, peerID, userID, fromID uint64, seq uint32, text string) {
	peerName := strconv.FormatUint(peerID, 10)
	userName := strconv.FormatUint(userID, 10)
	fromName := strconv.FormatUint(fromID, 10)
	if peerID == 0 {
		if contact, ok := c.contacts[userID]; ok {
			userName = contact.Remark
		}
		if contact, ok := c.contacts[fromID]; ok {
			fromName = contact.Remark
		}
	} else {
		if channel, ok := c.channels[peerID]; ok {
			peerName = channel.GroupName
		}
		if cmember, ok := c.cmembers[peerID][userID]; ok {
			userName = cmember.AutoRemark
		}
		if cmember, ok := c.cmembers[peerID][fromID]; ok {
			fromName = cmember.AutoRemark
		}
	}
	if log.GetLevel() > zerolog.DebugLevel {
		if peerID == 0 {
			log.PrintMessageSimple(time, userName, fromName, seq, text)
		} else if userID == 0 {
			log.PrintMessageSimple(time, peerName, fromName, seq, text)
		} else {
			log.PrintMessageSimple(time, userName+"("+peerName+")", fromName, seq, text)
		}
	} else {
		log.PrintMessage(time, peerName, userName, fromName, peerID, userID, fromID, seq, text)
	}
}

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
