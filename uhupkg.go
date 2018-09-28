package main

import (
	"bytes"

	_ "github.com/updatehub/updatehub/installmodes/copy"
	_ "github.com/updatehub/updatehub/installmodes/flash"
	_ "github.com/updatehub/updatehub/installmodes/imxkobs"
	_ "github.com/updatehub/updatehub/installmodes/raw"
	_ "github.com/updatehub/updatehub/installmodes/tarball"
	_ "github.com/updatehub/updatehub/installmodes/ubifs"

	"github.com/updatehub/updatehub/libarchive"
	"github.com/updatehub/updatehub/metadata"
)

func parsePackage(file string) (*metadata.UpdateMetadata, []byte, []byte, error) {
	la := &libarchive.LibArchive{}

	reader, err := libarchive.NewReader(la, file, 10240)
	if err != nil {
		return nil, nil, nil, err
	}
	defer reader.Free()

	data := bytes.NewBuffer(nil)
	err = reader.ExtractFile("metadata", data)
	if err != nil {
		return nil, nil, nil, err
	}

	metadata, err := metadata.NewUpdateMetadata(data.Bytes())
	if err != nil {
		return nil, nil, nil, err
	}

	reader, err = libarchive.NewReader(la, file, 10240)
	if err != nil {
		return metadata, data.Bytes(), nil, err
	}

	signature := bytes.NewBuffer(nil)
	err = reader.ExtractFile("signature", signature)
	if err != nil {
		return metadata, data.Bytes(), nil, err
	}

	return metadata, data.Bytes(), signature.Bytes(), nil
}
