package client

import (
	"strconv"

	"github.com/elap5e/go-mobileqq-api/mobileqq/client/db"
)

func (c *Client) dbInsertMessageSequence(uin uint64, v *db.MessageSequence) error {
	table := "u" + strconv.FormatUint(uin, 10) + "_message_sequences"
	tx, err := c.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	res, err := db.UpdateMessageSequenceTx(table, tx, v)
	if err != nil {
		return err
	}
	aff, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if aff == 0 {
		_, err := db.InsertMessageSequenceTx(table, tx, v)
		if err != nil {
			return err
		}
	}
	return tx.Commit()
}

func (c *Client) dbSelectMessageSequence(uin uint64, v *db.MessageSequence) error {
	table := "u" + strconv.FormatUint(uin, 10) + "_message_sequences"
	tx, err := c.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	if err := db.SelectMessageSequenceTx(table, tx, v); err != nil {
		return err
	}
	return tx.Commit()
}

func (c *Client) dbUpdateMessageSequence(uin uint64, v *db.MessageSequence) error {
	table := "u" + strconv.FormatUint(uin, 10) + "_message_sequences"
	tx, err := c.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	_, err = db.UpdateMessageSequenceTx(table, tx, v)
	if err != nil {
		return err
	}
	return tx.Commit()
}
