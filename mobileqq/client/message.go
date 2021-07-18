package client

import (
	"fmt"
	"strconv"

	"github.com/rs/zerolog"

	"github.com/elap5e/go-mobileqq-api/log"
	"github.com/elap5e/go-mobileqq-api/mobileqq/client/db"
	"github.com/elap5e/go-mobileqq-api/pb"
)

func (c *Client) PrintMessageRecord(mr *db.MessageRecord) {
	peerName := strconv.FormatInt(mr.PeerID, 10)
	userName := strconv.FormatInt(mr.UserID, 10)
	fromName := strconv.FormatInt(mr.FromID, 10)
	if mr.PeerID == 0 {
		if contact, ok := c.contacts[mr.UserID]; ok {
			userName = contact.Remark
		}
		if contact, ok := c.contacts[mr.FromID]; ok {
			fromName = contact.Remark
		}
	} else {
		if channel, ok := c.channels[mr.PeerID]; ok {
			peerName = channel.Name
		}
		if cmember, ok := c.cmembers[mr.PeerID][mr.UserID]; ok {
			userName = cmember.Remark
			if userName == "" {
				userName = cmember.Nick
			}
		}
		if cmember, ok := c.cmembers[mr.PeerID][mr.FromID]; ok {
			fromName = cmember.Remark
			if fromName == "" {
				fromName = cmember.Nick
			}
		}
	}
	if log.GetLevel() > zerolog.DebugLevel {
		if mr.PeerID == 0 {
			log.PrintMessageSimple(mr.Time, userName, fromName, mr.Seq, mr.Text)
		} else if mr.UserID == 0 {
			log.PrintMessageSimple(mr.Time, peerName, fromName, mr.Seq, mr.Text)
		} else {
			log.PrintMessageSimple(mr.Time, userName+"("+peerName+")", fromName, mr.Seq, mr.Text)
		}
	} else {
		log.PrintMessage(mr.Time, peerName, userName, fromName, mr.PeerID, mr.UserID, mr.FromID, mr.Seq, mr.Text)
	}
}

func syncUinPairMessage(uinPairMessage *pb.UinPairMessage) {
	log.Info().
		Int64("@peer", uinPairMessage.GetPeerUin()).
		Int64("readAt", uinPairMessage.GetLastReadTime()).
		Msgf("<-> [sync] %d message(s)", len(uinPairMessage.GetMessages()))

	for _, msg := range uinPairMessage.GetMessages() {
		peerID := msg.GetMessageHead().GetC2CTempMessageHead().GetGroupCode()
		userID := uinPairMessage.GetPeerUin()
		if msg.GetMessageHead().GetC2CCmd() == 0 {
			peerID = uinPairMessage.GetPeerUin()
			userID = 0
		}

		log.Debug().
			Str("@chat", fmt.Sprintf("@%du%d", peerID, userID)).
			Int32("@seq", msg.GetMessageHead().GetMessageSeq()).
			Int64("from", msg.GetMessageHead().GetFromUin()).
			Int64("time", msg.GetMessageHead().GetMessageTime()).
			Int32("type", msg.GetMessageHead().GetMessageType()).
			Int64("uid", msg.GetMessageHead().GetMessageUid()).
			Msg("--> [recv]")
	}
}
