/*
 * UpdateHub
 * Copyright (C) 2017
 * O.S. Systems Sofware LTDA: contato@ossystems.com.br
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package zephyr

import (
	"testing"

	"github.com/UpdateHub/updatehub-ce/installmodes"
	"github.com/stretchr/testify/assert"
)

func TestZephyrInit(t *testing.T) {
	val, err := installmodes.GetObject("zephyr")
	assert.NoError(t, err)

	r1, ok := val.(*ZephyrObject)
	if !ok {
		t.Error("Failed to cast return value of \"installmodes.GetObject()\" to ZephyrObject")
	}

	r2, ok := getObject().(*ZephyrObject)
	if !ok {
		t.Error("Failed to cast return value of \"getObject()\" to ZephyrObject")
	}

	assert.Equal(t, r2, r1)
}
