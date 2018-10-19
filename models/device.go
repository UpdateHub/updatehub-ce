package models

import (
	"time"

	"github.com/asdine/storm"
)

type Device struct {
	UID              string            `storm:"id" json:"uid"`
	Hardware         string            `json:"hardware"`
	Version          string            `json:"version"`
	ProductUID       string            `json:"product_uid"`
	DeviceIdentity   map[string]string `json:"device_identity"`
	DeviceAttributes map[string]string `json:"device_attributes"`
	Status           string            `json:"status"`
	LastSeen         time.Time         `json:"last_seen"`
}

func (d *Device) ActiveRollout(db *storm.DB) (*Rollout, error) {
	var rollouts []Rollout
	if err := db.Find("Running", true, &rollouts); err != nil && err != storm.ErrNotFound {
		return nil, err
	}

	var rollout *Rollout

	for _, r := range rollouts {
		for _, uid := range r.Devices {
			if uid == d.UID {
				rollout = &r
				break
			}
		}

		if rollout != nil {
			break
		}
	}

	return rollout, nil
}
