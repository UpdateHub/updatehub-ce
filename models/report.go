package models

import (
	"time"
)

type Report struct {
	ID        int       `storm:"id,increment" json:"id"`
	Device    string    `storm:"index" json:"device"`
	Rollout   int       `storm:"index" json:"rollout"`
	Status    string    `json:"status"`
	IsError   bool      `json:"error"`
	Message   string    `json:"message"`
	Timestamp time.Time `json:"timestamp"`
}
