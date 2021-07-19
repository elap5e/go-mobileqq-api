package client

import (
	"strconv"

	"github.com/elap5e/go-mobileqq-api/mobileqq/client/db"
)

func (c *Client) dbInsertMessageRecord(uin uint64, v *db.MessageRecord) error {
	table := "u" + strconv.FormatUint(uin, 10) + "_message_records"
	tx, err := c.db.Begin()
	if err != nil {
		return err
	}
	res, err := db.UpdateMessageRecordTx(table, tx, v)
	if err != nil {
		return err
	}
	aff, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if aff == 0 {
		_, err := db.InsertMessageRecordTx(table, tx, v)
		if err != nil {
			return err
		}
	}
	return tx.Commit()
}

func (c *Client) dbUpdateMessageRecord(uin uint64, v *db.MessageRecord) error {
	table := "u" + strconv.FormatUint(uin, 10) + "_message_records"
	tx, err := c.db.Begin()
	if err != nil {
		return err
	}
	_, err = db.UpdateMessageRecordTx(table, tx, v)
	if err != nil {
		return err
	}
	return tx.Commit()
}
