package migrations

import (
	"database/sql"
	"github.com/Vitaly-Baidin/l0/pkg/logging/zaplog"
	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigration(upAddUpdateAt, downAddUpdateAt)
}

func upAddUpdateAt(tx *sql.Tx) error {
	logger := zaplog.Logger
	_, err := tx.Exec(queryUp)
	if err != nil {
		logger.Errorf("failed work migrations (update): %v,\n", err)
		return err
	}
	return nil
}

func downAddUpdateAt(tx *sql.Tx) error {
	logger := zaplog.Logger
	_, err := tx.Exec(queryDown)
	if err != nil {
		logger.Errorf("failed work migrations (down): %v,\n", err)
		return err
	}
	return nil
}

const queryUp = `CREATE TABLE orders
(
    order_uid          VARCHAR(25) PRIMARY KEY,
    track_number       VARCHAR(50) UNIQUE NOT NULL,
    entry              VARCHAR(50)        NOT NULL,
    delivery           JSONB,
    payment            JSONB,
    items              JSONB,
    locale             VARCHAR(50)        NOT NULL,
    internal_signature VARCHAR(50),
    customer_id        VARCHAR(50)        NOT NULL,
    delivery_service   VARCHAR(50)        NOT NULL,
    shardkey           VARCHAR(50) UNIQUE NOT NULL,
    sm_id              INTEGER     UNIQUE NOT NULL,
    date_created       TIMESTAMP          NOT NULL,
    oof_shard          VARCHAR(50) UNIQUE NOT NULL
);
CREATE TABLE caches
(
    key        VARCHAR(250) UNIQUE NOT NULL PRIMARY KEY,
    value      jsonb,
    expiration int
);`

const queryDown = `DROP TABLE orders; DROP TABLE caches`
