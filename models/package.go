package models

type Package struct {
	UID               string   `storm:"id" json:"uid"`
	Version           string   `json:"version"`
	SupportedHardware []string `json:"supported_hardware"`
	Signature         []byte   `json:"signature"`
	Metadata          []byte   `json:"metadata"`
}
