package db

import (
	"database/sql"
)

func InsertChannelTx(table string, tx *sql.Tx, v *Channel) (sql.Result, error) {
	query := `INSERT INTO "` + table + `" ( "uin", "seq", "name", "memo", "member_num", "member_seq" )
VALUES( ?, ?, ?, ?, ?, ? );`
	stmt, err := tx.Prepare(query)
	if err != nil {
		return nil, err
	}
	return stmt.Exec(v.Uin, v.Seq, v.Name, v.Memo, v.MemberNum, v.MemberSeq)
}

func UpdateChannelTx(table string, tx *sql.Tx, v *Channel) (sql.Result, error) {
	query := `UPDATE "` + table + `" SET "seq" = ?, "name" = ?, "memo" = ?, "member_num" = ?, "member_seq" = ? WHERE "uin" = ? ;`
	stmt, err := tx.Prepare(query)
	if err != nil {
		return nil, err
	}
	return stmt.Exec(v.Seq, v.Name, v.Memo, v.MemberNum, v.MemberSeq, v.Uin)
}
