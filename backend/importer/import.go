package importer

import (
	"archive/tar"
	"compress/gzip"
	"errors"
	"io"
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
