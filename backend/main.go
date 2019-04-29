package main

import (
	"database/sql"
	"os"

	"./data"
	"./importer"
)

const dbFilename = "db.sqlite3"

// URL is represented as unicode character codes to hinder other candidates cheating by Googling the URL.
const url = "\u0068\u0074\u0074\u0070\u0073\u003a\u002f\u002f\u0073\u0033\u002e\u0061\u006d\u0061\u007a\u006f\u006e\u0061\u0077\u0073\u002e\u0063\u006f\u006d\u002f\u0072\u006d\u002d\u0072\u0061\u006e\u0074\u002d\u0069\u006e\u0074\u0065\u0072\u0076\u0069\u0065\u0077\u0069\u006e\u0067\u002f\u0070\u0072\u006f\u0064\u0075\u0063\u0074\u0073\u002e\u0074\u0061\u0072\u002e\u0067\u007a"

func main() {
	if _, err := os.Stat(dbFilename); err == nil {
		if err := os.Remove(dbFilename); err != nil {
			panic(err)
		}
	}

	data.WithTx(dbFilename, func(tx *sql.Tx) {
		adRepo := data.NewTxAdvertiserRepository(tx)
		prodRepo := data.NewTxProductRepository(tx)
		persister := importer.NewPersister(adRepo, prodRepo)

		err := importer.ImportOTF(url, persister)
		if err != nil {
			panic(err)
		}

		tx.Commit()
	})

	println("Done!")
}
