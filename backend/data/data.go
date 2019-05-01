package data

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

const schemaSql = `
CREATE TABLE advertiser (
  id integer PRIMARY KEY AUTOINCREMENT,
  name text NOT NULL UNIQUE
);

CREATE TABLE product (
  sku nvarchar(36) PRIMARY KEY,
  name text NOT NULL UNIQUE,
  advertiser_id int NOT NULL,

  FOREIGN KEY (advertiser_id) REFERENCES advertiser (id)
);
`

func withTxTest(action func(tx *sql.Tx)) {
	WithTx(":memory:", action)
}

func WithTx(source string, action func(tx *sql.Tx)) {
	db, err := sql.Open("sqlite3", source)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	db.Exec(schemaSql)

	tx, err := db.Begin()
	if err != nil {
		panic(err)
	}

	action(tx)
}
