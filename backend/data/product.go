package data

import (
	"database/sql"
	"errors"

	"github.com/mattn/go-sqlite3"
	_ "github.com/mattn/go-sqlite3"
)

type ProductEntity struct {
	Sku          string
	Name         string
	AdvertiserId int
}

type ProductRepository interface {
	Save(ProductEntity) error
}

type TxSqlProductRepository struct {
	tx *sql.Tx
}

func NewTxProductRepository(tx *sql.Tx) ProductRepository {
	return TxSqlProductRepository{tx}
}

func (repo TxSqlProductRepository) Save(item ProductEntity) error {
	_, err := repo.tx.Exec("INSERT INTO product (sku, name, advertiser_id) VALUES (?, ?, ?)", item.Sku, item.Name, item.AdvertiserId)

	if err, ok := err.(sqlite3.Error); ok {
		if err.Code == sqlite3.ErrConstraint {
			return ErrDuplicateProduct
		}
	}

	return err
}

var ErrDuplicateProduct = errors.New("Duplicate product")
