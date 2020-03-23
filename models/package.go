// Copyright (C) 2018 O.S. Systems Sofware LTDA
//
// SPDX-License-Identifier: MIT

package models

type Package struct {
	UID               string      `storm:"id" json:"uid"`
	Version           string      `json:"version"`
	SupportedHardware interface{} `json:"supported_hardware"`
	Signature         []byte      `json:"signature"`
	Metadata          []byte      `json:"metadata"`
}
