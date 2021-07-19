package client

import (
	"strconv"
)

func (c *Client) dbCreateAccountTable() error {
	queries := []string{`CREATE TABLE IF NOT EXISTS "accounts" (
  "id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
  "uin" INTEGER NOT NULL,
  "sync_cookie" BLOB NULL,
  "created_at" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "deleted_at" TIMESTAMP NULL,
  "updated_at" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  UNIQUE (
    "uin" ASC,
    "updated_at" ASC
  )
);`,
		`CREATE INDEX IF NOT EXISTS "accounts_idx"
ON "accounts" (
  "id" ASC,
  "uin" ASC
);`,
		`CREATE TRIGGER IF NOT EXISTS "accounts_updated_at"
AFTER UPDATE
ON "accounts"
FOR EACH ROW
BEGIN
UPDATE "accounts" SET "updated_at" = CURRENT_TIMESTAMP WHERE id = old.id;
END;`}
	for _, query := range queries {
		stmt, err := c.db.Prepare(query)
		if err != nil {
			return err
		}
		_, err = stmt.Exec()
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *Client) dbCreateChannelTableByUin(uin uint64) error {
	channelMemberTable := "u" + strconv.FormatUint(uin, 10) + "_channel_members"
	queries := []string{`CREATE TABLE IF NOT EXISTS "` + channelMemberTable + `" (
  "id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
  "time" INTEGER NOT NULL,
  "uin" INTEGER NOT NULL,
  "nick" TEXT NOT NULL,
  "gender" INTEGER NOT NULL,
  "remark" TEXT NOT NULL,
  "status" INTEGER NOT NULL,
  "channel_id" INTEGER NOT NULL,
  "created_at" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "deleted_at" TIMESTAMP NULL,
  "updated_at" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  UNIQUE (
    "time" ASC,
    "uin" ASC,
    "channel_id" ASC
  ),
  CONSTRAINT "fk_channel_id" FOREIGN KEY (
    "channel_id"
  ) REFERENCES ` + channelMemberTable + ` (
    "uin"
  ) ON DELETE SET DEFAULT ON UPDATE NO ACTION
);`,
		`CREATE INDEX IF NOT EXISTS "` + channelMemberTable + `_channel_id_idx"
ON "` + channelMemberTable + `" (
  "channel_id" ASC
);`,
		`CREATE INDEX IF NOT EXISTS "` + channelMemberTable + `_idx"
ON "` + channelMemberTable + `" (
  "id" ASC,
  "time" ASC
);`,
		`CREATE INDEX IF NOT EXISTS "` + channelMemberTable + `_uin_idx"
ON "` + channelMemberTable + `" (
  "uin" ASC
);`,
		`CREATE TRIGGER IF NOT EXISTS "` + channelMemberTable + `_updated_at"
AFTER UPDATE
ON "` + channelMemberTable + `"
FOR EACH ROW
BEGIN
UPDATE "` + channelMemberTable + `" SET "updated_at" = CURRENT_TIMESTAMP WHERE id = old.id;
END;`}
	channelTable := "u" + strconv.FormatUint(uin, 10) + "_channels"
	queries = append(queries, []string{`CREATE TABLE IF NOT EXISTS "` + channelTable + `" (
  "id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
  "uin" INTEGER NOT NULL,
  "seq" INTEGER NOT NULL,
  "name" TEXT NOT NULL,
  "memo" TEXT NOT NULL,
  "member_num" INTEGER NOT NULL,
  "member_seq" INTEGER NOT NULL,
  "created_at" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "deleted_at" TIMESTAMP NULL,
  "updated_at" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  UNIQUE (
    "uin" ASC
  )
);`,
		`CREATE INDEX IF NOT EXISTS ` + channelTable + `_idx
ON "` + channelTable + `" (
  "id" ASC,
  "uin" ASC
);`,
		`CREATE TRIGGER IF NOT EXISTS "` + channelTable + `_updated_at"
AFTER UPDATE
ON "` + channelTable + `"
FOR EACH ROW
BEGIN
UPDATE "` + channelTable + `" SET "updated_at" = CURRENT_TIMESTAMP WHERE id = old.id;
END;`}...)
	for _, query := range queries {
		stmt, err := c.db.Prepare(query)
		if err != nil {
			return err
		}
		_, err = stmt.Exec()
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *Client) dbCreateContactTableByUin(uin uint64) error {
	contactGroupTable := "u" + strconv.FormatUint(uin, 10) + "_contact_groups"
	queries := []string{`CREATE TABLE IF NOT EXISTS "` + contactGroupTable + `" (
  "id" INTEGER NOT NULL PRIMARY KEY,
  "name" TEXT NOT NULL,
  "updated_at" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);`,
		`CREATE TRIGGER IF NOT EXISTS "` + contactGroupTable + `_updated_at"
AFTER UPDATE
ON "` + contactGroupTable + `"
FOR EACH ROW
BEGIN
UPDATE "` + contactGroupTable + `" SET "updated_at" = CURRENT_TIMESTAMP WHERE id = old.id;
END;`}
	contactTable := "u" + strconv.FormatUint(uin, 10) + "_contacts"
	queries = append(queries, []string{`CREATE TABLE IF NOT EXISTS "` + contactTable + `" (
  "id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
  "uin" INTEGER NOT NULL,
  "nick" TEXT NOT NULL,
  "gender" INTEGER NOT NULL,
  "remark" TEXT NOT NULL,
  "status" INTEGER NOT NULL,
  "group_id" INTEGER NOT NULL,
  "created_at" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "deleted_at" TIMESTAMP NULL,
  "updated_at" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  UNIQUE (
    "uin" ASC
  ),
  CONSTRAINT "fk_group_id" FOREIGN KEY (
    "group_id"
  ) REFERENCES ` + contactGroupTable + ` (
    "id"
  ) ON DELETE SET DEFAULT ON UPDATE NO ACTION
);`,
		`CREATE INDEX IF NOT EXISTS ` + contactTable + `_group_id_idx
ON "` + contactTable + `" (
  "group_id" ASC
);`,
		`CREATE INDEX IF NOT EXISTS ` + contactTable + `_idx
ON "` + contactTable + `" (
  "id" ASC,
  "uin" ASC
);`,
		`CREATE TRIGGER IF NOT EXISTS "` + contactTable + `_updated_at"
AFTER UPDATE
ON "` + contactTable + `"
FOR EACH ROW
BEGIN
UPDATE "` + contactTable + `" SET "updated_at" = CURRENT_TIMESTAMP WHERE id = old.id;
END;`}...)
	for _, query := range queries {
		stmt, err := c.db.Prepare(query)
		if err != nil {
			return err
		}
		_, err = stmt.Exec()
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *Client) dbCreateMessageRecordTableByUin(uin uint64) error {
	table := "u" + strconv.FormatUint(uin, 10) + "_message_records"
	queries := []string{`CREATE TABLE IF NOT EXISTS "` + table + `" (
  "id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
  "time" INTEGER NOT NULL,
  "seq" INTEGER NOT NULL,
  "uid" INTEGER NOT NULL,
  "peer_id" INTEGER NOT NULL,
  "user_id" INTEGER NOT NULL,
  "from_id" INTEGER NOT NULL,
  "text" TEXT NOT NULL,
  "type" INTEGER NOT NULL,
  "created_at" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "deleted_at" TIMESTAMP NULL,
  "updated_at" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  UNIQUE (
    "time" ASC,
    "seq" ASC,
    "peer_id" ASC,
    "user_id" ASC,
    "text" ASC
  )
);`,
		`CREATE INDEX IF NOT EXISTS "` + table + `_from_id_idx"
ON "` + table + `" (
  "from_id" ASC
);`,
		`CREATE INDEX IF NOT EXISTS "` + table + `_idx"
ON "` + table + `" (
  "id" ASC,
  "time" ASC
);`,
		`CREATE INDEX IF NOT EXISTS "` + table + `_peer_id_idx"
ON "` + table + `" (
  "peer_id" ASC
);`,
		`CREATE INDEX IF NOT EXISTS "` + table + `_seq_idx"
ON "` + table + `" (
  "seq" ASC
);`,
		`CREATE INDEX IF NOT EXISTS "` + table + `_uid_idx"
ON "` + table + `" (
  "uid" ASC
);`,
		`CREATE INDEX IF NOT EXISTS "` + table + `_user_id_idx"
ON "` + table + `" (
  "user_id" ASC
);`,
		`CREATE TRIGGER IF NOT EXISTS "` + table + `_updated_at"
AFTER UPDATE
ON "` + table + `"
FOR EACH ROW
BEGIN
UPDATE "` + table + `" SET "updated_at" = CURRENT_TIMESTAMP WHERE id = old.id;
END;`}
	for _, query := range queries {
		stmt, err := c.db.Prepare(query)
		if err != nil {
			return err
		}
		_, err = stmt.Exec()
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *Client) dbCreateMessageSequenceTableByUin(uin uint64) error {
	table := "u" + strconv.FormatUint(uin, 10) + "_message_sequences"
	queries := []string{`CREATE TABLE IF NOT EXISTS "` + table + `" (
  "id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
  "peer_id" INTEGER NOT NULL,
  "user_id" INTEGER NOT NULL,
  "type" INTEGER NOT NULL,
  "max_seq" INTEGER NOT NULL,
  "created_at" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "deleted_at" TIMESTAMP NULL,
  "updated_at" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  UNIQUE (
    "peer_id" ASC,
    "user_id" ASC,
    "type" ASC
  )
);`,
		`CREATE TRIGGER IF NOT EXISTS "` + table + `_updated_at"
AFTER UPDATE
ON "` + table + `"
FOR EACH ROW
BEGIN
UPDATE "` + table + `" SET "updated_at" = CURRENT_TIMESTAMP WHERE id = old.id;
END;`}
	for _, query := range queries {
		stmt, err := c.db.Prepare(query)
		if err != nil {
			return err
		}
		_, err = stmt.Exec()
		if err != nil {
			return err
		}
	}
	return nil
}
