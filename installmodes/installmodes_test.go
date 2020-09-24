/*
 * UpdateHub
 * Copyright (C) 2017
 * O.S. Systems Sofware LTDA: contato@ossystems.com.br
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package installmodes

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

type TestObject struct {
}

func TestRegisterInstallMode(t *testing.T) {
	mode := RegisterInstallMode(InstallMode{
		Name:      "test",
		GetObject: func() interface{} { return &TestObject{} },
	})

	defer mode.Unregister()

	obj, err := GetObject("test")
	if !assert.NotNil(t, obj) {
		t.Fatal(err)
	}
}

func TestUnregisterInstallMode(t *testing.T) {
	mode := RegisterInstallMode(InstallMode{
		Name:      "test",
		GetObject: func() interface{} { return &TestObject{} },
	})

	defer mode.Unregister()

	obj, err := GetObject("test")
	if !assert.NotNil(t, obj) {
		t.Fatal(err)
	}

	mode.Unregister()

	obj, err = GetObject("test")
	if !assert.Nil(t, obj) {
		t.Fatal(err)
	}
}

func TestGetObjectNotFound(t *testing.T) {
	_, err := GetObject("test")

	if assert.Error(t, err) {
		assert.Equal(t, errors.New("Object not found"), err)
	}
}
