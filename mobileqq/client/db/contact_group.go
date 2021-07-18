package db

import (
	"database/sql"
)

func InsertContactGroupTx(table string, tx *sql.Tx, v *ContactGroup) (sql.Result, error) {
	query := `INSERT INTO "` + table + `" ( "id", "name" )
VALUES( ?, ? );`
	stmt, err := tx.Prepare(query)
	if err != nil {
		return nil, err
	}
	return stmt.Exec(v.ID, v.Name)
}

func UpdateContactGroupTx(table string, tx *sql.Tx, v *ContactGroup) (sql.Result, error) {
	query := `UPDATE "` + table + `" SET "name" = ? WHERE "id" = ? ;`
	stmt, err := tx.Prepare(query)
	if err != nil {
		return nil, err
	}
	return stmt.Exec(v.Name, v.ID)
}
