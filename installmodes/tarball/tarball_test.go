/*
 * UpdateHub
 * Copyright (C) 2017
 * O.S. Systems Sofware LTDA: contato@ossystems.com.br
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package tarball

import (
	"testing"

	"github.com/UpdateHub/updatehub-ce/installmodes"
	"github.com/stretchr/testify/assert"
)

func TestTarballInit(t *testing.T) {
	val, err := installmodes.GetObject("tarball")
	assert.NoError(t, err)

	tb1, ok := val.(*TarballObject)
	if !ok {
		t.Error("Failed to cast return value of \"installmodes.GetObject()\" to TarballObject")
	}

	tb2, ok := getObject().(*TarballObject)
	if !ok {
		t.Error("Failed to cast return value of \"getObject()\" to TarballObject")
	}

	assert.Equal(t, tb2, tb1)
}
