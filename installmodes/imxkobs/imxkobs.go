/*
 * UpdateHub
 * Copyright (C) 2017
 * O.S. Systems Sofware LTDA: contato@ossystems.com.br
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package imxkobs

import (
	"github.com/UpdateHub/updatehub-ce/installmodes"
	"github.com/UpdateHub/updatehub-ce/metadata"
)

func init() {
	installmodes.RegisterInstallMode(installmodes.InstallMode{
		Name:      "imxkobs",
		GetObject: getObject,
	})
}

func getObject() interface{} {
	return &ImxKobsObject{}
}

// ImxKobsObject encapsulates the "imxkobs" handler data and functions
type ImxKobsObject struct {
	metadata.ObjectMetadata

	Add1KPadding    bool   `json:"1k_padding,omitempty"`
	SearchExponent  int    `json:"search_exponent,omitempty"`
	Chip0DevicePath string `json:"chip_0_device_path,omitempty"`
	Chip1DevicePath string `json:"chip_1_device_path,omitempty"`
}
