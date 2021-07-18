package db

import (
	"database/sql"
)

func InsertContactTx(table string, tx *sql.Tx, v *Contact) (sql.Result, error) {
	query := `INSERT INTO "` + table + `" ( "uin", "nick", "gender", "remark", "status", "group_id" )
VALUES( ?, ?, ?, ?, ?, ? );`
	stmt, err := tx.Prepare(query)
	if err != nil {
		return nil, err
	}
	return stmt.Exec(v.Uin, v.Nick, v.Gender, v.Remark, v.Status, v.GroupID)
}

func UpdateContactTx(table string, tx *sql.Tx, v *Contact) (sql.Result, error) {
	query := `UPDATE "` + table + `" SET "nick" = ?, "gender" = ?, "remark" = ?, "status" = ?, "group_id" = ? WHERE "uin" = ? ;`
	stmt, err := tx.Prepare(query)
	if err != nil {
		return nil, err
	}
	return stmt.Exec(v.Nick, v.Gender, v.Remark, v.Status, v.GroupID, v.Uin)
}
