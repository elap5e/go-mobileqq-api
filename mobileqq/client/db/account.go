package db

import (
	"database/sql"
)

func InsertAccountTx(tx *sql.Tx, v *Account) (sql.Result, error) {
	query := `INSERT INTO "accounts" ( "uin", "sync_cookie" )
VALUES( ?, ? );`
	stmt, err := tx.Prepare(query)
	if err != nil {
		return nil, err
	}
	return stmt.Exec(v.Uin, v.SyncCookie)
}

func SelectAccountTx(tx *sql.Tx, uin int64) (*Account, error) {
	query := `SELECT "id", "uin", "sync_cookie" FROM "accounts" WHERE "uin" = ? AND "deleted_at" IS NULL ;`
	row := tx.QueryRow(query, uin)
	if err := row.Err(); err != nil {
		return nil, err
	}
	v := Account{}
	if err := row.Scan(&v.ID, &v.Uin, &v.SyncCookie); err != nil {
		return nil, err
	}
	return &v, nil
}

func UpdateAccountTx(tx *sql.Tx, v *Account) (sql.Result, error) {
	update, err := SelectAccountTx(tx, v.Uin)
	if err != nil {
		return nil, err
	}
	if len(v.SyncCookie) != 0 {
		update.SyncCookie = v.SyncCookie
	}
	query := `UPDATE "accounts" SET "sync_cookie" = ? WHERE "uin" = ? AND "deleted_at" IS NULL ;`
	stmt, err := tx.Prepare(query)
	if err != nil {
		return nil, err
	}
	return stmt.Exec(update.SyncCookie, update.Uin)
}
