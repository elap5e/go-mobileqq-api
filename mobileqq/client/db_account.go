package client

import (
	"database/sql"

	"github.com/elap5e/go-mobileqq-api/mobileqq/client/db"
)

func (c *Client) dbInsertAccount(v *db.Account) error {
	tx, err := c.db.Begin()
	if err != nil {
		return err
	}
	res, err := db.UpdateAccountTx(tx, v)
	if err != nil && err != sql.ErrNoRows {
		return err
	}
	aff := int64(0)
	if res != nil {
		aff, err = res.RowsAffected()
		if err != nil {
			return err
		}
	}
	if aff == 0 {
		_, err := db.InsertAccountTx(tx, v)
		if err != nil {
			return err
		}
	}
	return tx.Commit()
}

func (c *Client) dbSelectAccount(uin uint64) (*db.Account, error) {
	tx, err := c.db.Begin()
	if err != nil {
		return nil, err
	}
	v, err := db.SelectAccountTx(tx, int64(uin))
	if err != nil {
		return nil, err
	}
	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return v, nil
}

func (c *Client) dbUpdateAccount(v *db.Account) error {
	tx, err := c.db.Begin()
	if err != nil {
		return err
	}
	_, err = db.UpdateAccountTx(tx, v)
	if err != nil {
		return err
	}
	return tx.Commit()
}
