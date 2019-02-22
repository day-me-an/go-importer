package importer

import (
	"io"
	"io/ioutil"
	"os"
	"strings"
	"sync"
	"testing"
)

func TestDecompressArchive(t *testing.T) {
	var err error

	r, err := os.Open("fixtures/test.tar.gz")
	if err != nil {
		t.Fatal("Couldn't open test fixture file")
	}
	defer r.Close()

	err = DecompressArchive(r, func(filename string, entryReader io.Reader) {
		readContent := func() string {
			b, _ := ioutil.ReadAll(entryReader)
			return string(b)
		}

		switch filename {
		case "hello.txt":
			if content := readContent(); content != "world" {
				t.Errorf("Expected 'world' but found '%s'", content)
			}

		case "lol.txt":
			if content := readContent(); content != "haha" {
				t.Errorf("Expected 'haha' but found '%s'", content)
			}

		default:
			t.Errorf("Unexpected file found in test archive: %s", filename)
		}
	})

	if err != nil {
		t.Errorf("Decompressing archive failed due to %s", err.Error())
	}
}

func TestExtractAds_Empty(t *testing.T) {
	extractAdvertisersHelper(t, "", []Advertiser{})
}

func TestExtractAds_Single(t *testing.T) {
	extractAdvertisersHelper(t, "nike", []Advertiser{{Name: "nike"}})
}

func TestExtractAds_TrailingComma(t *testing.T) {
	extractAdvertisersHelper(t, "nike,", []Advertiser{{Name: "nike"}})
}

func TestExtractAds_Quoted(t *testing.T) {
	extractAdvertisersHelper(t, "\"nike\"", []Advertiser{{Name: "nike"}})
}

func TestExtractAds_Multiple(t *testing.T) {
	extractAdvertisersHelper(t, "nike,apple", []Advertiser{
		{Name: "nike"},
		{Name: "apple"},
	})
}

func TestExtractAds_MultipleSpaced(t *testing.T) {
	extractAdvertisersHelper(t, " nike, apple ,google, ", []Advertiser{
		{Name: "nike"},
		{Name: "apple"},
		{Name: "google"},
	})
}

func extractAdvertisersHelper(t *testing.T, content string, expected []Advertiser) {
	output := make(chan Advertiser, len(expected))
	defer close(output)

	r := strings.NewReader(content)
	go ExtractAdvertisers(r, output)

	for i := 0; i < len(expected); i++ {
		// TODO: handle nothing being there with some kind of timeout.
		found := <-output

		if found != expected[i] {
			t.Errorf("Expected %s but found %s at %d", expected[i], found, i)
		}
	}

	select {
	case extra := <-output:
		t.Errorf("It's returning extra data: %s", extra)
	default:
		// Noop because it's as expected.
	}
}

func TestExtractProducts_Empty(t *testing.T) {
	extractProductsHelper(t, "", []Product{})
}

func TestExtractProducts_HeaderOnly(t *testing.T) {
	extractProductsHelper(t, "name,sku,advertiser\n", []Product{})
}

func TestExtractProducts_Data(t *testing.T) {
	extractProductsHelper(t, "name,sku,advertiser\niphone,123,google", []Product{
		Product{Name: "iphone", Sku: "123", Advertiser: "google"},
	})
}

func TestExtractProducts_Multiple(t *testing.T) {
	extractProductsHelper(t, "name,sku,advertiser\niphone,123,google\nmacbook,456,facebook", []Product{
		Product{Name: "iphone", Sku: "123", Advertiser: "google"},
		Product{Name: "macbook", Sku: "456", Advertiser: "facebook"},
	})
}

func TestExtractProducts_WrongNumberOfFields(t *testing.T) {
	// It should just skip these records.
	extractProductsHelper(t, "name,sku,advertiser\niphone,123,google,EXTRA_DATA", []Product{})
}

/*
Unfortunately, Go lacks generics, so some code will have to be duplicated to maintain static typing.
Empty interfaces were a possibility, but the resulting code seemed overcomplicated due to slices.
Reflection was another option, but that would obviously be much worse.

It appears Golang V2 will have generics.
*/
func extractProductsHelper(t *testing.T, content string, expected []Product) {
	output := make(chan Product, len(expected))

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()

		// Will terminate when channel is closed.
		i := 0
		for found := range output {
			if i < len(expected) {
				if found != expected[i] {
					t.Errorf("Expected %s but found %s at %d", expected[i], found, i)
				}
			} else {
				t.Errorf("It's returning extra data: %s", found)
			}
			i++
		}
	}()

	r := strings.NewReader(content)
	err := ExtractProducts(r, output)
	close(output)

	if err != nil {
		t.Error("Extracting failed due to", err)
	}

	wg.Wait()
}
