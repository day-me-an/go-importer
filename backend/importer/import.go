package importer

import (
	"archive/tar"
	"bufio"
	"compress/gzip"
	"errors"
	"io"
	"strings"
)

func DecompressArchive(r io.Reader, onFile func(filename string, r io.Reader)) error {
	decompressed, err := gzip.NewReader(r)
	if err != nil {
		return errors.New("Couldn't create gzip reader")
	}

	tarReader := tar.NewReader(decompressed)
	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return errors.New("Couldn't read next tar archive entry")
		}

		// We only care about files.
		if header.Typeflag == tar.TypeReg {
			onFile(header.Name, tarReader)
		}
	}

	return nil
}

type Advertiser struct {
	Name string
}

func ExtractAdvertisers(r io.Reader, output chan Advertiser) {
	scanner := bufio.NewScanner(r)
	scanner.Split(cslSplitter)

	for scanner.Scan() {
		name := scanner.Text()
		name = strings.Trim(name, " ")
		name = strings.Replace(name, "\"", "", -1)

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
