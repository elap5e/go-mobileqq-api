package client

import (
	"strconv"

	"github.com/elap5e/go-mobileqq-api/mobileqq/client/db"
)

func (c *Client) dbInsertChannel(uin uint64, v *db.Channel) error {
	table := "u" + strconv.FormatUint(uin, 10) + "_channels"
	tx, err := c.db.Begin()
	if err != nil {
		return err
	}
	res, err := db.UpdateChannelTx(table, tx, v)
	if err != nil {
		return err
	}
	aff, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if aff == 0 {
		_, err := db.InsertChannelTx(table, tx, v)
		if err != nil {
			return err
		}
	}
	return tx.Commit()
}

func (c *Client) dbUpdateChannel(uin uint64, v *db.Channel) error {
	table := "u" + strconv.FormatUint(uin, 10) + "_channels"
	tx, err := c.db.Begin()
	if err != nil {
		return err
	}
	_, err = db.UpdateChannelTx(table, tx, v)
	if err != nil {
		return err
	}
	return tx.Commit()
}
