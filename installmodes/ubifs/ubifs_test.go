/*
 * UpdateHub
 * Copyright (C) 2017
 * O.S. Systems Sofware LTDA: contato@ossystems.com.br
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package ubifs

import (
	"testing"

	"github.com/UpdateHub/updatehub-ce/installmodes"
	"github.com/stretchr/testify/assert"
)

func TestUbifsInit(t *testing.T) {
	val, err := installmodes.GetObject("ubifs")
	assert.NoError(t, err)

	f1, ok := val.(*UbifsObject)
	if !ok {
		t.Error("Failed to cast return value of \"installmodes.GetObject()\" to UbifsObject")
	}

	f2, ok := getObject().(*UbifsObject)
	if !ok {
		t.Error("Failed to cast return value of \"getObject()\" to UbifsObject")
	}

	assert.Equal(t, f2, f1)
}
