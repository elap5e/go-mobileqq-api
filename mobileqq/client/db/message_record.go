package db

import (
	"database/sql"
)

func DeleteMessageRecordTx(table string, tx *sql.Tx, v *MessageRecord) (sql.Result, error) {
	query := `UPDATE "` + table + `" SET "deleted_at" = ? WHERE "time" >= ? AND "seq" = ? AND "peer_id" = ? AND "user_id" = ? ;`
	stmt, err := tx.Prepare(query)
	if err != nil {
		return nil, err
	}
	return stmt.Exec(v.DeletedAt, v.Time, v.Seq, v.PeerID, v.UserID)
}

func InsertMessageRecordTx(table string, tx *sql.Tx, v *MessageRecord) (sql.Result, error) {
	query := `INSERT INTO "` + table + `" ( "time", "seq", "uid", "peer_id", "user_id", "from_id", "text", "type" )
VALUES( ?, ?, ?, ?, ?, ?, ?, ? );`
	stmt, err := tx.Prepare(query)
	if err != nil {
		return nil, err
	}
	return stmt.Exec(v.Time, v.Seq, v.Uid, v.PeerID, v.UserID, v.FromID, v.Text, v.Type)
}

func UpdateMessageRecordTx(table string, tx *sql.Tx, v *MessageRecord) (sql.Result, error) {
	query := `UPDATE "` + table + `" SET "uid" = ?, "from_id" = ?, "text" = ? WHERE "time" >= ? AND "seq" = ? AND "peer_id" = ? AND "user_id" = ? AND "type" = ? ;`
	stmt, err := tx.Prepare(query)
	if err != nil {
		return nil, err
	}
	return stmt.Exec(v.Uid, v.FromID, v.Text, v.Time, v.Seq, v.PeerID, v.UserID, v.Type)
}
