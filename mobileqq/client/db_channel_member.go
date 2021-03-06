package client

import (
	"strconv"

	"github.com/elap5e/go-mobileqq-api/mobileqq/client/db"
)

func (c *Client) dbInsertChannelMember(uin uint64, v *db.ChannelMember) error {
	table := "u" + strconv.FormatUint(uin, 10) + "_channel_members"
	tx, err := c.db.Begin()
	if err != nil {
		return err
	}
	res, err := db.UpdateChannelMemberTx(table, tx, v)
	if err != nil {
		return err
	}
	aff, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if aff == 0 {
		_, err := db.InsertChannelMemberTx(table, tx, v)
		if err != nil {
			return err
		}
	}
	return tx.Commit()
}

func (c *Client) dbUpdateChannelMember(uin uint64, v *db.ChannelMember) error {
	table := "u" + strconv.FormatUint(uin, 10) + "_channel_members"
	tx, err := c.db.Begin()
	if err != nil {
		return err
	}
	_, err = db.UpdateChannelMemberTx(table, tx, v)
	if err != nil {
		return err
	}
	return tx.Commit()
}
