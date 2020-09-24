/*
 * UpdateHub
 * Copyright (C) 2017
 * O.S. Systems Sofware LTDA: contato@ossystems.com.br
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package copy

import (
	"github.com/UpdateHub/updatehub-ce/installmodes"
	"github.com/UpdateHub/updatehub-ce/metadata"
)

func init() {
	installmodes.RegisterInstallMode(installmodes.InstallMode{
		Name:      "copy",
		GetObject: getObject,
	})
}

func getObject() interface{} {
	return &CopyObject{
		ChunkSize: 128 * 1024,
	}
}

// CopyObject encapsulates the "copy" handler data and functions
type CopyObject struct {
	metadata.ObjectMetadata
	metadata.CompressedObject

	Target        string      `json:"target"`
	TargetType    string      `json:"target-type"`
	TargetPath    string      `json:"target-path"`
	TargetGID     interface{} `json:"target-gid"` // can be string or int
	TargetUID     interface{} `json:"target-uid"` // can be string or int
	TargetMode    string      `json:"target-mode"`
	FSType        string      `json:"filesystem"`
	FormatOptions string      `json:"format-options,omitempty"`
	MustFormat    bool        `json:"format?,omitempty"`
	MountOptions  string      `json:"mount-options,omitempty"`
	ChunkSize     int         `json:"chunk-size,omitempty"`
}
