package db

import (
	"database/sql"
)

func InsertChannelMemberTx(table string, tx *sql.Tx, v *ChannelMember) (sql.Result, error) {
	query := `INSERT INTO "` + table + `" ( "time", "uin", "nick", "gender", "remark", "status", "channel_id" )
VALUES( ?, ?, ?, ?, ?, ?, ? );`
	stmt, err := tx.Prepare(query)
	if err != nil {
		return nil, err
	}
	return stmt.Exec(v.Time, v.Uin, v.Nick, v.Gender, v.Remark, v.Status, v.ChannelID)
}

func UpdateChannelMemberTx(table string, tx *sql.Tx, v *ChannelMember) (sql.Result, error) {
	query := `UPDATE "` + table + `" SET "nick" = ?, "gender" = ?, "remark" = ?, "status" = ? WHERE "time" = ? AND "uin" = ? AND "channel_id" = ? ;`
	stmt, err := tx.Prepare(query)
	if err != nil {
		return nil, err
	}
	return stmt.Exec(v.Nick, v.Gender, v.Remark, v.Status, v.Time, v.Uin, v.ChannelID)
}
