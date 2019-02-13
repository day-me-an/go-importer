package importer

import (
	"io"
	"io/ioutil"
	"os"
	"strings"
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
	output := make(chan Advertiser)
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
