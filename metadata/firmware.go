/*
 * UpdateHub
 * Copyright (C) 2017
 * O.S. Systems Sofware LTDA: contato@ossystems.com.br
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package metadata

type FirmwareMetadata struct {
	ProductUID       string            `json:"product-uid"`
	DeviceIdentity   map[string]string `json:"device-identity"`
	Version          string            `json:"version"`
	Hardware         string            `json:"hardware"`
	DeviceAttributes map[string]string `json:"device-attributes"`
}
