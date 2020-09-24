/*
 * UpdateHub
 * Copyright (C) 2019
 * O.S. Systems Sofware LTDA: contato@ossystems.com.br
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package zephyr

import (
	"github.com/UpdateHub/updatehub-ce/installmodes"
	"github.com/UpdateHub/updatehub-ce/metadata"
)

func init() {
	installmodes.RegisterInstallMode(installmodes.InstallMode{
		Name:      "zephyr",
		GetObject: getObject,
	})
}

func getObject() interface{} {
	return &ZephyrObject{}
}

// ZephyrObject encapsulates the "zephyr" handler data and functions
type ZephyrObject struct {
	metadata.ObjectMetadata
}
