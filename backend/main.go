package main

import "./importer"
import "./data"

func main() {
	adRepo := data.NewTxAdvertiserRepository(nil)
	prodRepo := data.NewTxProductRepository(nil)
	persister := importer.NewPersister(adRepo, prodRepo)

	// URL is represented as unicode character codes to hinder other candidates cheating by Googling the URL.
	url := "\u0068\u0074\u0074\u0070\u0073\u003a\u002f\u002f\u0073\u0033\u002e\u0061\u006d\u0061\u007a\u006f\u006e\u0061\u0077\u0073\u002e\u0063\u006f\u006d\u002f\u0072\u006d\u002d\u0072\u0061\u006e\u0074\u002d\u0069\u006e\u0074\u0065\u0072\u0076\u0069\u0065\u0077\u0069\u006e\u0067\u002f\u0070\u0072\u006f\u0064\u0075\u0063\u0074\u0073\u002e\u0074\u0061\u0072\u002e\u0067\u007a"
	err := importer.ImportOTF(url, persister)
	if err != nil {
		panic(err)
	}

	println("Done!")
}
