/*
 * UpdateHub
 * Copyright (C) 2017
 * O.S. Systems Sofware LTDA: contato@ossystems.com.br
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package raw_delta

import (
	"testing"

	"github.com/UpdateHub/updatehub-ce/installmodes"
	"github.com/stretchr/testify/assert"
)

func TestRawDeltaInit(t *testing.T) {
	val, err := installmodes.GetObject("raw-delta")
	assert.NoError(t, err)

	r1, ok := val.(*RawDeltaObject)
	if !ok {
		t.Error("Failed to cast return value of \"installmodes.GetObject()\" to RawDeltaObject")
	}

	r2, ok := getObject().(*RawDeltaObject)
	if !ok {
		t.Error("Failed to cast return value of \"getObject()\" to RawDeltaObject")
	}

	assert.Equal(t, r2, r1)
}
