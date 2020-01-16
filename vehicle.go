package main

import "encoding/json"

// Vehicle is a vehicle.
type Vehicle struct {
	Linked         `pg:"-"`
	ManufacturerID string        `pg:",pk" json:"-"`
	Manufacturer   *Manufacturer `json:"-"`
	PowerSourceID  int           `json:"-"`
	PowerSource    *PowerSource  `json:"-"`
	TSN            string        `pg:"id,pk" json:"tsn,omitempty"`
	TradeName      string        `json:"tradeName,omitempty"`
	CommercialName string        `json:"commercialName,omitempty"`
	AllotmentDate  string        `json:"allotmentDate,omitempty"`
	Category       string        `json:"category,omitempty"`
	Bodywork       string        `json:"bodywork,omitempty"`
	Power          int           `json:"power,omitempty"`
	EngineCapacity int           `json:"engineCapacity,omitempty"`
	Axles          int           `json:"axles,omitempty"`
	PoweredAxles   int           `json:"poweredAxles,omitempty"`
	Seats          int           `json:"seats,omitempty"`
	MaximumMass    int           `json:"maximumMass,omitempty"`
}

func (v *Vehicle) String() string {
	bytes, _ := json.Marshal(v)
	return string(bytes)
}
