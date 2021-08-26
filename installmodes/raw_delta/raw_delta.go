/*
 * UpdateHub
 * Copyright (C) 2021
 * O.S. Systems Sofware LTDA: contato@ossystems.com.br
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package raw_delta

import (
	"github.com/UpdateHub/updatehub-ce/installmodes"
	"github.com/UpdateHub/updatehub-ce/metadata"
)

func init() {
	installmodes.RegisterInstallMode(installmodes.InstallMode{
		Name:      "raw-delta",
		GetObject: getObject,
	})
}

func getObject() interface{} {
	return &RawDeltaObject{}
}

// RawDeltaObject encapsulates the "raw-delta" handler data and functions
type RawDeltaObject struct {
	metadata.ObjectMetadata

	Target     string `json:"target"`
	TargetType string `json:"target-type"`
	ChunkSize  int    `json:"chunk-size,omitempty"`
	Seek       int    `json:"seek,omitempty"`
}
