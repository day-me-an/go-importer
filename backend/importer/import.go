package importer

import (
	"archive/tar"
	"bufio"
	"compress/gzip"
	"encoding/csv"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"sync/atomic"

	"../data"
)

const ProductConsumersNum = 4

type Advertiser struct {
	Name string
}

type Product struct {
	Sku        string
	Name       string
	Advertiser string
}

// Imports the records contained in the compressed archive file on-the-fly while it's being read without fully downloading it first.
func ImportOTF(url string, persister Persister) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	DecompressArchive(resp.Body, func(filename string, r io.Reader) {
		fmt.Printf("Processing %s\n", filename)

		var wg sync.WaitGroup
		var done, imported, skipped, duplicates, failed uint64
		var doneMutex sync.Mutex

		switch filename {
		case "advertisers.txt":
			queue := make(chan Advertiser)
			// Just start _one_ consumer for processing the advertisers since it's so small.
			wg.Add(1)
			go func() {
				defer wg.Done()

				for item := range queue {
					persister.SaveAdvertiser(item)
					atomic.AddUint64(&imported, 1)
				}
			}()
			ExtractAdvertisers(r, queue)
			close(queue)

		case "products.csv":
			queue := make(chan Product)
			// Start multiple consumers to process them in parallel.
			wg.Add(ProductConsumersNum)
			for i := 0; i < ProductConsumersNum; i++ {
				go func() {
					defer wg.Done()

					for item := range queue {
						if err := persister.SaveProduct(item); err == nil {
							atomic.AddUint64(&imported, 1)
						} else if err == ErrUnknownAdvertiser {
							atomic.AddUint64(&skipped, 1)
						} else if err == data.ErrDuplicateProduct {
							atomic.AddUint64(&duplicates, 1)
						} else {
							fmt.Println("Failed due to", err)
							atomic.AddUint64(&failed, 1)
						}

						// Display some indication of progress.
						doneMutex.Lock()
						if done%10000 == 0 {
							fmt.Printf("Done %d\n", done)
						}
						doneMutex.Unlock()

						atomic.AddUint64(&done, 1)
					}
				}()
			}
			err := ExtractProducts(r, queue)
			close(queue)
			if err != nil {
				panic(err)
			}
		}

		wg.Wait()
		fmt.Printf("Imported %d records, skipped %d, %d duplicates and failed on %d in %s\n", imported, skipped, duplicates, failed, filename)
	})

	return nil
}

func DecompressArchive(r io.Reader, onFile func(filename string, r io.Reader)) error {
	decompressed, err := gzip.NewReader(r)
	if err != nil {
		return err
	}

	tarReader := tar.NewReader(decompressed)
	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		// We only care about files.
		if header.Typeflag == tar.TypeReg {
			onFile(header.Name, tarReader)
		}
	}

	return nil
}

func ExtractAdvertisers(r io.Reader, output chan<- Advertiser) {
	scanner := bufio.NewScanner(r)
	scanner.Split(cslSplitter)

	for scanner.Scan() {
		name := scanner.Text()
		name = strings.Trim(name, " ")
		// All values are wrapped in single quotes.
		name = strings.Trim(name, "'")

		if len(name) > 0 {
			output <- Advertiser{Name: name}
		}
	}
}

func cslSplitter(data []byte, atEOF bool) (advance int, token []byte, err error) {
	for i := 0; i < len(data); i++ {
		if data[i] == ',' {
			return i + 1, data[:i], nil
		}
	}
	return 0, data, bufio.ErrFinalToken
}

func ExtractProducts(r io.Reader, output chan<- Product) error {
	csvReader := csv.NewReader(r)

	// Skip the header.
	csvReader.Read()

	for {
		record, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		// Just skip invalid records.
		if parseErr, ok := err.(*csv.ParseError); ok && parseErr.Err == csv.ErrFieldCount {
			continue
		}
		if err != nil {
			return err
		}

		output <- Product{
			Name:       record[0],
			Sku:        record[1],
			Advertiser: record[2],
		}
	}

	return nil
}
