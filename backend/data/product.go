package data

import "database/sql"

type ProductEntity struct {
	Sku          string
	Name         string
	AdvertiserId int
}

type ProductRepository interface {
	Save(ProductEntity)
}

type TxSqlProductRepository struct {
	tx *sql.Tx
}

func NewTxProductRepository(tx *sql.Tx) ProductRepository {
	return TxSqlProductRepository{tx}
}

func (repo TxSqlProductRepository) Save(item ProductEntity) {
	repo.tx.Exec("INSERT INTO product (sku, name, advertiser_id) VALUES (?, ?, ?)", item.Sku, item.Name, item.AdvertiserId)
}
