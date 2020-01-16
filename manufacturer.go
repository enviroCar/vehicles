package main

import (
	"encoding/json"
)

// Manufacturer is a vehicle manufacturer.
type Manufacturer struct {
	Linked `pg:"-"`
	ID     string `pg:",pk" json:"hsn,omitempty"`
	Name   string `json:"name,omitempty"`
}

func (m *Manufacturer) String() string {
	bytes, _ := json.Marshal(m)
	return string(bytes)
}
