/*
 * UpdateHub
 * Copyright (C) 2017
 * O.S. Systems Sofware LTDA: contato@ossystems.com.br
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package flash

import (
	"testing"

	"github.com/UpdateHub/updatehub-ce/installmodes"
	"github.com/stretchr/testify/assert"
)

func TestFlashInit(t *testing.T) {
	val, err := installmodes.GetObject("flash")
	assert.NoError(t, err)

	f1, ok := val.(*FlashObject)
	if !ok {
		t.Error("Failed to cast return value of \"installmodes.GetObject()\" to FlashObject")
	}

	f2, ok := getObject().(*FlashObject)
	if !ok {
		t.Error("Failed to cast return value of \"getObject()\" to FlashObject")
	}

	assert.Equal(t, f2, f1)
}
