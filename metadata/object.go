/*
 * UpdateHub
 * Copyright (C) 2017
 * O.S. Systems Sofware LTDA: contato@ossystems.com.br
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package metadata

import (
	"encoding/json"

	"github.com/UpdateHub/updatehub-ce/installmodes"
)

// ObjectMetadata contains the common properties of a package's object from JSON metadata
type ObjectMetadata struct {
	Object `json:"-"`

	Sha256sum          string      `json:"sha256sum"`
	Mode               string      `json:"mode"`
	Size               int64       `json:"size"`
	Compressed         bool        `json:"bool"`
	InstallIfDifferent interface{} `json:"install-if-different,omitempty"`
}

func NewObjectMetadata(bytes []byte) (Object, error) {
	var v map[string]interface{}

	err := json.Unmarshal(bytes, &v)
	if err != nil {
		return nil, err
	}

	var obj Object

	o, err := installmodes.GetObject(v["mode"].(string))
	if err == nil {
		obj = o.(Object)
	} else {
		return nil, err
	}

	json.Unmarshal(bytes, &obj)

	return obj, nil
}

func (o ObjectMetadata) GetObjectMetadata() ObjectMetadata {
	return o
}

type CompressedObject struct {
	CompressedSize   float64 `json:"required-compressed-size"`
	UncompressedSize float64 `json:"required-uncompressed-size"`
}

type Object interface {
	GetObjectMetadata() ObjectMetadata
}
