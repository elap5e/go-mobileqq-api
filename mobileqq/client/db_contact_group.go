package client

import (
	"strconv"

	"github.com/elap5e/go-mobileqq-api/mobileqq/client/db"
)

func (c *Client) dbInsertContactGroup(uin uint64, v *db.ContactGroup) error {
	table := "u" + strconv.FormatUint(uin, 10) + "_contact_group"
	tx, err := c.db.Begin()
	if err != nil {
		return err
	}
	res, err := db.UpdateContactGroupTx(table, tx, v)
	if err != nil {
		return err
	}
	aff, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if aff == 0 {
		_, err := db.InsertContactGroupTx(table, tx, v)
		if err != nil {
			return err
		}
	}
	return tx.Commit()
}

func (c *Client) dbUpdateContactGroup(uin uint64, v *db.ContactGroup) error {
	table := "u" + strconv.FormatUint(uin, 10) + "_contact_group"
	tx, err := c.db.Begin()
	if err != nil {
		return err
	}
	_, err = db.UpdateContactGroupTx(table, tx, v)
	if err != nil {
		return err
	}
	return tx.Commit()
}
