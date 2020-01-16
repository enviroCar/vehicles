package main

import (
	"encoding/json"
)

// PowerSource is the power source of a vehicle.
type PowerSource struct {
	Linked      `pg:"-"`
	ID          int    `pg:",pk" json:"id,omitempty"`
	ShortName   string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
}

func (ps *PowerSource) String() string {
	bytes, _ := json.Marshal(ps)
	return string(bytes)
}
