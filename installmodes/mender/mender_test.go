/*
 * UpdateHub
 * Copyright (C) 2017
 * O.S. Systems Sofware LTDA: contato@ossystems.com.br
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package mender

import (
	"testing"

	"github.com/UpdateHub/updatehub-ce/installmodes"
	"github.com/stretchr/testify/assert"
)

func TestMenderInit(t *testing.T) {
	val, err := installmodes.GetObject("mender")
	assert.NoError(t, err)

	r1, ok := val.(*MenderObject)
	if !ok {
		t.Error("Failed to cast return value of \"installmodes.GetObject()\" to MenderObject")
	}

	r2 := &MenderObject{}

	assert.Equal(t, r2, r1)
}
