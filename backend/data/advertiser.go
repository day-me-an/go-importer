package data

import "database/sql"

type AdvertiserEntity struct {
	Id   int
	Name string
}

type AdvertiserRepository interface {
	Save(AdvertiserEntity)
	GetByName(string) *AdvertiserEntity
}

type TxSqlAdvertiserRepository struct {
	/*
		TODO: research whether the database/sql package has an interface that contains the common
		querying function signatures, which would allow support using this repository outside of a transaction.
	*/
	tx *sql.Tx
}

func NewTxAdvertiserRepository(tx *sql.Tx) AdvertiserRepository {
	return TxSqlAdvertiserRepository{tx}
}

func (repo TxSqlAdvertiserRepository) Save(item AdvertiserEntity) {
	// TODO: add an 'ON CONFLICT' to update it if it already exists.
	repo.tx.Exec("INSERT INTO advertiser (name) VALUES (?)", item.Name)
}

func (repo TxSqlAdvertiserRepository) GetByName(name string) *AdvertiserEntity {
	row := repo.tx.QueryRow("SELECT id FROM advertiser WHERE name = ?", name)

	entity := &AdvertiserEntity{Name: name}
	if err := row.Scan(&entity.Id); err == sql.ErrNoRows {
		return nil
	} else if err != nil {
		panic(err)
	}
	return entity
}
