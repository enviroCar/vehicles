package main

import (
	"io"
	"strconv"

	"github.com/go-pg/pg/v9"
)

// Repository is the vehicle repository.
type Repository struct{ db *pg.DB }

// NewRepository creates a new Repository.
func NewRepository(options *pg.Options) *Repository {
	return &Repository{db: pg.Connect(options)}
}

var _ io.Closer = (*Repository)(nil)

// Close closes this repository.
func (r *Repository) Close() error {
	return r.db.Close()
}

// GetManufacturers returns all manufacturers.
func (r *Repository) GetManufacturers() ([]*Manufacturer, error) {
	var entities []*Manufacturer
	err := r.db.Model(&entities).Select()
	if err != nil {
		if err == pg.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return entities, err
}

// GetManufacturer returns the specified manufacturer.
func (r *Repository) GetManufacturer(id string) (*Manufacturer, error) {
	manufacturer := new(Manufacturer)
	err := r.db.Model(manufacturer).Where("id = ? ", id).First()
	if err != nil {
		if err == pg.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return manufacturer, nil
}

// GetVehicles returns all vehicles of the manufacturer.
func (r *Repository) GetVehicles(manufacturer *Manufacturer) ([]*Vehicle, error) {
	var vehicles []*Vehicle
	err := r.db.Model(&vehicles).
		Column("id", "trade_name", "commercial_name", "allotment_date", "manufacturer_id").
		Where("manufacturer_id = ?", manufacturer.ID).
		Select()
	if err != nil {
		if err == pg.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return vehicles, nil
}

// GetVehicle tries to get the specified vehicle.
func (r *Repository) GetVehicle(manufacturer *Manufacturer, id string) (*Vehicle, error) {

	vehicle := new(Vehicle)
	err := r.db.Model(vehicle).
		Relation("Manufacturer").
		Relation("PowerSource").
		Where("vehicle.manufacturer_id = ? AND vehicle.id = ?", manufacturer.ID, id).
		First()
	if err != nil {
		if err == pg.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return vehicle, nil
}

// GetPowerSources gets all available power sources.
func (r *Repository) GetPowerSources() ([]*PowerSource, error) {
	var entities []*PowerSource
	err := r.db.Model(&entities).Select()
	if err != nil {
		if err == pg.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return entities, nil
}

// GetPowerSource gets the specified power source.
func (r *Repository) GetPowerSource(id string) (*PowerSource, error) {
	nid, err := strconv.ParseInt(id, 10, 32)
	if err != nil {
		return nil, NewErrNotFoundF("id is bad '%v'", id)
	}

	powerSource := new(PowerSource)
	err = r.db.Model(powerSource).
		Where("id = ? ", nid).
		First()
	if err != nil {
		if err == pg.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return powerSource, nil
}
