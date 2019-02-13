package importer

import (
	"io"
	"io/ioutil"
	"os"
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
