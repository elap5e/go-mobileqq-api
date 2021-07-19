package db

import (
	"database/sql"
)

func InsertMessageSequenceTx(table string, tx *sql.Tx, v *MessageSequence) (sql.Result, error) {
	query := `INSERT INTO "` + table + `" ( "peer_id", "user_id", "type", "max_seq" )
VALUES( ?, ?, ?, ? );`
	stmt, err := tx.Prepare(query)
	if err != nil {
		return nil, err
	}
	return stmt.Exec(v.PeerID, v.UserID, v.Type, v.MaxSeq)
}

func SelectMessageSequenceTx(table string, tx *sql.Tx, v *MessageSequence) error {
	query := `SELECT "max_seq" FROM "` + table + `" WHERE "peer_id" = ? AND "user_id" = ? AND "type" = ? ;`
	row := tx.QueryRow(query, v.PeerID, v.UserID, v.Type)
	if err := row.Err(); err != nil {
		return err
	}
	if err := row.Scan(&v.MaxSeq); err != nil {
		return err
	}
	return nil
}

func UpdateMessageSequenceTx(table string, tx *sql.Tx, v *MessageSequence) (sql.Result, error) {
	query := `UPDATE "` + table + `" SET "max_seq" = ? WHERE "peer_id" = ? AND "user_id" = ? AND "type" = ? ;`
	stmt, err := tx.Prepare(query)
	if err != nil {
		return nil, err
	}
	return stmt.Exec(v.MaxSeq, v.PeerID, v.UserID, v.Type)
}
