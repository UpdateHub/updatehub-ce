/*
 * UpdateHub
 * Copyright (C) 2017
 * O.S. Systems Sofware LTDA: contato@ossystems.com.br
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package raw

import (
	"testing"

	"github.com/UpdateHub/updatehub-ce/installmodes"
	"github.com/stretchr/testify/assert"
)

func TestRawInit(t *testing.T) {
	val, err := installmodes.GetObject("raw")
	assert.NoError(t, err)

	r1, ok := val.(*RawObject)
	if !ok {
		t.Error("Failed to cast return value of \"installmodes.GetObject()\" to RawObject")
	}

	r2, ok := getObject().(*RawObject)
	if !ok {
		t.Error("Failed to cast return value of \"getObject()\" to RawObject")
	}

	assert.Equal(t, r2, r1)
}
