package data

import (
	"database/sql"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func TestAdvertiser_GetByName_Empty(t *testing.T) {
	withTx(func(tx *sql.Tx) {
		repo := TxSqlAdvertiserRepository{tx: tx}
		actual := repo.GetByName("google")

		if actual != nil {
			t.Error("Expected nothing, but got", actual)
		}
	})
}

func TestAdvertiser_GetByName_Some(t *testing.T) {
	withTx(func(tx *sql.Tx) {
		tx.Exec("INSERT INTO advertiser (name) VALUES ('google')")

		repo := TxSqlAdvertiserRepository{tx: tx}
		actual := repo.GetByName("google")

		if actual == nil {
			t.Error("Expected something, but got nothing")
		}
		if actual.Id != 1 || actual.Name != "google" {
			t.Error("Got wrong entity", actual)
		}
	})
}

func TestAdvertiser_Save(t *testing.T) {
	withTx(func(tx *sql.Tx) {
		repo := TxSqlAdvertiserRepository{tx: tx}
		repo.Save(AdvertiserEntity{Id: -1, Name: "google"})

		row := tx.QueryRow("SELECT * FROM advertiser")
		var actual AdvertiserEntity
		row.Scan(&actual.Id, &actual.Name)

		if actual.Id != 1 || actual.Name != "google" {
			t.Error("Saved wrong entity", actual)
		}
	})
}

func withTx(action func(tx *sql.Tx)) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	db.Exec(`
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
`)

	tx, err := db.Begin()
	if err != nil {
		panic(err)
	}

	action(tx)
}
