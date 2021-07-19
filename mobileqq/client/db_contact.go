package client

import (
	"strconv"

	"github.com/elap5e/go-mobileqq-api/mobileqq/client/db"
)

func (c *Client) dbInsertContact(uin uint64, v *db.Contact) error {
	table := "u" + strconv.FormatUint(uin, 10) + "_contacts"
	tx, err := c.db.Begin()
	if err != nil {
		return err
	}
	res, err := db.UpdateContactTx(table, tx, v)
	if err != nil {
		return err
	}
	aff, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if aff == 0 {
		_, err := db.InsertContactTx(table, tx, v)
		if err != nil {
			return err
		}
	}
	return tx.Commit()
}

func (c *Client) dbUpdateContact(uin uint64, v *db.Contact) error {
	table := "u" + strconv.FormatUint(uin, 10) + "_contacts"
	tx, err := c.db.Begin()
	if err != nil {
		return err
	}
	_, err = db.UpdateContactTx(table, tx, v)
	if err != nil {
		return err
	}
	return tx.Commit()
}
