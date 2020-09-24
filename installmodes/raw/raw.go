/*
 * UpdateHub
 * Copyright (C) 2017
 * O.S. Systems Sofware LTDA: contato@ossystems.com.br
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package raw

import (
	"github.com/UpdateHub/updatehub-ce/installmodes"
	"github.com/UpdateHub/updatehub-ce/metadata"
)

func init() {
	installmodes.RegisterInstallMode(installmodes.InstallMode{
		Name:      "raw",
		GetObject: getObject,
	})
}

func getObject() interface{} {
	return &RawObject{
		ChunkSize: 128 * 1024,
		Count:     -1,
		Truncate:  true,
	}
}

// RawObject encapsulates the "raw" handler data and functions
type RawObject struct {
	metadata.ObjectMetadata
	metadata.CompressedObject

	Target     string `json:"target"`
	TargetType string `json:"target-type"`
	ChunkSize  int    `json:"chunk-size,omitempty"`
	Skip       int    `json:"skip,omitempty"`
	Seek       int    `json:"seek,omitempty"`
	Count      int    `json:"count,omitempty"`
	Truncate   bool   `json:"truncate,omitempty"`
}
