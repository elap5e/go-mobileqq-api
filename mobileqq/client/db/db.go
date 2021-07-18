package db

import (
	"time"

	"github.com/elap5e/go-mobileqq-api/pb"
)

type Channel struct {
	ID        int64
	Uin       int64
	Seq       int32
	Name      string
	Memo      string
	MemberNum int32
	MemberSeq int32
	CreatedAt time.Time
	DeletedAt time.Time
	UpdatedAt time.Time
}

type ChannelMember struct {
	ID        int64
	Time      int64
	Uin       int64
	Nick      string
	Gender    int8
	Remark    string
	Status    int8
	ChannelID int64
	CreatedAt time.Time
	DeletedAt time.Time
	UpdatedAt time.Time
}

type Contact struct {
	ID        int64
	Uin       int64
	Nick      string
	Gender    int8
	Remark    string
	Status    int8
	GroupID   int8
	CreatedAt time.Time
	DeletedAt time.Time
	UpdatedAt time.Time
}

type ContactGroup struct {
	ID        int8
	Name      string
	UpdatedAt time.Time
}

type MessageRecord struct {
	ID        int64
	Time      int64
	Seq       int32
	Uid       int64
	PeerID    int64
	UserID    int64
	FromID    int64
	Text      string
	Elements  []*pb.Element
	Type      int32
	CreatedAt time.Time
	DeletedAt time.Time
	UpdatedAt time.Time
}
