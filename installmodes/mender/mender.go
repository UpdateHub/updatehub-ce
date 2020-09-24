/*
 * UpdateHub
 * Copyright (C) 2018
 * O.S. Systems Sofware LTDA: contato@ossystems.com.br
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package mender

import (
	"github.com/UpdateHub/updatehub-ce/installmodes"
	"github.com/UpdateHub/updatehub-ce/metadata"
)

func init() {
	installmodes.RegisterInstallMode(installmodes.InstallMode{
		Name:      "mender",
		GetObject: getObject,
	})
}

func getObject() interface{} {
	return &MenderObject{}
}

// MenderObject encapsulates the "mender" handler data and functions
type MenderObject struct {
	metadata.ObjectMetadata
}
