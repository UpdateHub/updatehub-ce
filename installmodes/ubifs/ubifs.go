/*
 * UpdateHub
 * Copyright (C) 2017
 * O.S. Systems Sofware LTDA: contato@ossystems.com.br
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package ubifs

import (
	"github.com/UpdateHub/updatehub-ce/installmodes"
	"github.com/UpdateHub/updatehub-ce/metadata"
)

func init() {
	installmodes.RegisterInstallMode(installmodes.InstallMode{
		Name:      "ubifs",
		GetObject: getObject,
	})
}

func getObject() interface{} {
	return &UbifsObject{}
}

// UbifsObject encapsulates the "ubifs" handler data and functions
type UbifsObject struct {
	metadata.ObjectMetadata
	metadata.CompressedObject

	Target     string `json:"target"`
	TargetType string `json:"target-type"`
}
