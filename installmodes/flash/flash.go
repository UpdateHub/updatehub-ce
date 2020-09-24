/*
 * UpdateHub
 * Copyright (C) 2017
 * O.S. Systems Sofware LTDA: contato@ossystems.com.br
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package flash

import (
	"github.com/UpdateHub/updatehub-ce/installmodes"
	"github.com/UpdateHub/updatehub-ce/metadata"
)

func init() {
	installmodes.RegisterInstallMode(installmodes.InstallMode{
		Name:      "flash",
		GetObject: getObject,
	})
}

func getObject() interface{} {
	return &FlashObject{}
}

// FlashObject encapsulates the "flash" handler data and functions
type FlashObject struct {
	metadata.ObjectMetadata

	Target     string `json:"target"`
	TargetType string `json:"target-type"`
}
