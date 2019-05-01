package data

import (
	"database/sql"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func TestAdvertiser_GetByName_Empty(t *testing.T) {
	withTxTest(func(tx *sql.Tx) {
		repo := TxSqlAdvertiserRepository{tx: tx}
		actual := repo.GetByName("google")

		if actual != nil {
			t.Error("Expected nothing, but got", actual)
		}
	})
}

func TestAdvertiser_GetByName_Some(t *testing.T) {
	withTxTest(func(tx *sql.Tx) {
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
	withTxTest(func(tx *sql.Tx) {
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
