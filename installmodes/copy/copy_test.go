/*
 * UpdateHub
 * Copyright (C) 2017
 * O.S. Systems Sofware LTDA: contato@ossystems.com.br
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package copy

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/UpdateHub/updatehub-ce/installmodes"
)

func TestCopyInit(t *testing.T) {
	val, err := installmodes.GetObject("copy")
	assert.NoError(t, err)

	cp1, ok := val.(*CopyObject)
	if !ok {
		t.Error("Failed to cast return value of \"installmodes.GetObject()\" to CopyObject")
	}

	cp2, ok := getObject().(*CopyObject)
	if !ok {
		t.Error("Failed to cast return value of \"getObject()\" to CopyObject")
	}

	assert.Equal(t, cp2, cp1)
}
