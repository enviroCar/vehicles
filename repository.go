package main

import (
	"io"
	"strconv"

	"github.com/go-pg/pg/v9"
	"github.com/go-pg/pg/v9/orm"
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
	return entities, err
}

// GetManufacturer returns the specified manufacturer.
func (r *Repository) GetManufacturer(id string) (*Manufacturer, error) {
	var manufacturers []*Manufacturer
	err := r.db.Model(&manufacturers).Where("id = ? ", id).Limit(1).Select()
	if err != nil {
		return nil, err
	}
	if len(manufacturers) == 0 {
		return nil, ErrNotFound
	}
	return manufacturers[0], nil
}

// GetVehicles returns all vehicles of the manufacturer.
func (r *Repository) GetVehicles(manufacturer *Manufacturer) ([]*Vehicle, error) {
	var vehicles []*Vehicle
	err := r.db.Model(&vehicles).
		Column("id", "trade_name", "commercial_name", "allotment_date", "manufacturer_id").
		Where("manufacturer_id = ?", manufacturer.ID).
		Select()
	return vehicles, err
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
	return entities, err
}

// GetPowerSource gets the specified power source.
func (r *Repository) GetPowerSource(id string) (*PowerSource, error) {
	nid, err := strconv.ParseInt(id, 10, 32)
	if err != nil {
		return nil, NewErrBadRequestF("id is bad '%v'", id)
	}

	var powerSources []*PowerSource
	var q *orm.Query = r.db.Model(&powerSources)
	err = q.Where("id = ? ", nid).Limit(1).Select()
	if err != nil {
		return nil, err
	}
	if len(powerSources) == 0 {
		return nil, ErrNotFound
	}
	return powerSources[0], nil
}
