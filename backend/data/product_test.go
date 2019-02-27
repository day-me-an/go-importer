package data

import (
	"database/sql"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func TestProduct_Save(t *testing.T) {
	withTx(func(tx *sql.Tx) {
		repo := TxSqlProductRepository{tx: tx}
		repo.Save(ProductEntity{Sku: "123abc", Name: "iphone", AdvertiserId: 1})

		row := tx.QueryRow("SELECT * FROM product")
		var actual ProductEntity
		row.Scan(&actual.Sku, &actual.Name, &actual.AdvertiserId)

		if actual.Sku != "123abc" || actual.Name != "iphone" || actual.AdvertiserId != 1 {
			t.Error("Saved wrong entity", actual)
		}
	})
}
