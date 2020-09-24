/*
 * UpdateHub
 * Copyright (C) 2017
 * O.S. Systems Sofware LTDA: contato@ossystems.com.br
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package imxkobs

import (
	"testing"

	"github.com/UpdateHub/updatehub-ce/installmodes"
	"github.com/stretchr/testify/assert"
)

func TestImxKobsInit(t *testing.T) {
	val, err := installmodes.GetObject("imxkobs")
	assert.NoError(t, err)

	ik1, ok := val.(*ImxKobsObject)
	if !ok {
		t.Error("Failed to cast return value of \"installmodes.GetObject()\" to ImxKobsObject")
	}

	ik2, ok := getObject().(*ImxKobsObject)
	if !ok {
		t.Error("Failed to cast return value of \"getObject()\" to ImxKobsObject")
	}

	assert.Equal(t, ik2, ik1)
}
