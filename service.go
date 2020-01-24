package main

import (
	"io"
	"strconv"
)

// Service is the vehicle service.
type Service struct{ repository *Repository }

var _ io.Closer = (*Service)(nil)

// NewService creates a new Service
func NewService(repository *Repository) *Service {
	return &Service{repository}
}

// Close closes this service.
func (s *Service) Close() error {
	return s.repository.Close()
}

// GetRoot returns the root content.
func (s *Service) GetRoot(context *Context) (interface{}, error) {

	context.logger.Info("get root")

	links := &Linked{}
	href, err := context.URL(s.GetManufacturers)()
	if err != nil {
		context.logger.WithError(err).Error("could not create manufacturer links")
		return nil, ErrInternalServer
	}
	links.AddLink(NewLink(href, "manufacturers", "application/json", "Manufacturers"))

	href, err = context.URL(s.GetPowerSources)()
	if err != nil {
		context.logger.WithError(err).Error("could not create power source links")
		return nil, ErrInternalServer
	}
	links.AddLink(NewLink(href, "powerSources", "application/json", "Power Sources"))
	return links, nil
}

func (s *Service) manufacturerLink(context *Context, m *Manufacturer, relation string) (*Link, error) {
	href, err := context.URL(s.GetManufacturer)("hsn", m.ID)
	return NewLink(href, relation, "application/json", m.Name), err
}

func (s *Service) powerSourceLink(context *Context, m *PowerSource, relation string) (*Link, error) {
	href, err := context.URL(s.GetPowerSource)("id", strconv.Itoa(m.ID))
	return NewLink(href, relation, "application/json", m.ShortName), err
}

func (s *Service) vehicleLink(context *Context, vehicle *Vehicle, relation string) (*Link, error) {
	href, err := context.URL(s.GetVehicle)("hsn", vehicle.ManufacturerID, "tsn", vehicle.TSN)
	return NewLink(href, relation, "application/json", vehicle.CommercialName), err
}

// GetManufacturers returns all manufacturers.
func (s *Service) GetManufacturers(context *Context) (interface{}, error) {

	context.logger.Info("get manufacturers")

	entities, err := s.repository.GetManufacturers()
	if err != nil {
		if context.server.IsCriticalError(err)  {
			context.logger.WithError(err).Error("could not get manufacturers")
			return nil, ErrInternalServer
		}
		return nil, err
	}

	for _, m := range entities {
		link, err := s.manufacturerLink(context, m, "canonical")
		if err != nil {
			context.logger.WithError(err).Error("could not create manufacturer link")
			return nil, ErrInternalServer
		}
		m.AddLink(link)
	}

	return entities, nil
}

// GetManufacturer returns the specified manufacturer.
func (s *Service) GetManufacturer(context *Context) (interface{}, error) {

	hsn := context.Params["hsn"]

	context.logger.Infof("get manufacturer by id: '%s'", hsn)

	m, err := s.repository.GetManufacturer(hsn)
	if err != nil {
		if context.server.IsCriticalError(err)  {
			context.logger.WithError(err).Error("could not get manufacturer")
			return nil, ErrInternalServer
		}
		return nil, err
	}

	href, err := context.URL(s.GetVehicles)("hsn", m.ID)
	if err != nil {
		context.logger.WithError(err).Error("could not create vehicles links")
		return nil, ErrInternalServer
	}
	m.AddLink(NewLink(href, "vehicles", "application/json", ""))

	link, err := s.manufacturerLink(context, m, "self")
	if err != nil {
		context.logger.WithError(err).Error("could not create manufacturer self link")
		return nil, ErrInternalServer
	}
	m.AddLink(link)

	return m, nil
}

// GetVehicles returns all vehicles of the manufacturer.
func (s *Service) GetVehicles(context *Context) (interface{}, error) {

	context.logger.Infof("get vehicles")

	hsn := context.Params["hsn"]

	m, err := s.repository.GetManufacturer(hsn)
	if err != nil {
		if context.server.IsCriticalError(err) {
			context.logger.WithError(err).Errorf("could not get manufacturer by id: '%s'", hsn)
			return nil, ErrInternalServer
		}
		return nil, err
	}

	vehicles, err := s.repository.GetVehicles(m)
	if err != nil {
		if context.server.IsCriticalError(err) {
			context.logger.WithError(err).Errorf("could not get vehicles by manufacturer: %v", m)
			return nil, ErrInternalServer
		}
		return nil, err
	}

	for _, vehicle := range vehicles {
		link, err := s.vehicleLink(context, vehicle, "canonical")
		if err != nil {
			context.logger.WithError(err).Error("could not create vehicle link")
			return nil, ErrInternalServer
		}
		vehicle.AddLink(link)
	}
	return vehicles, err
}

// GetVehicle tries to get the specified vehicle.
func (s *Service) GetVehicle(context *Context) (interface{}, error) {

	context.logger.Infof("get vehicle by id and manufacturer")

	hsn := context.Params["hsn"]

	m, err := s.repository.GetManufacturer(hsn)
	if err != nil {
		if context.server.IsCriticalError(err) {
			context.logger.WithError(err).Errorf("could not get manufacturer by id: '%s'", hsn)
			return nil, ErrInternalServer
		}
		return nil, err
	}

	tsn := context.Params["tsn"]

	v, err := s.repository.GetVehicle(m, tsn)
	if err != nil {
		if context.server.IsCriticalError(err) {
			context.logger.WithError(err).Errorf("could not get vehicle by id: '%s'", tsn)
			return nil, ErrInternalServer
		}
		return nil, err
	}

	link, err := s.vehicleLink(context, v, "self")
	if err != nil {
		context.logger.WithError(err).Error("could not create vehicle link self")
		return nil, ErrInternalServer
	}
	v.AddLink(link)

	link, err = s.powerSourceLink(context, v.PowerSource, "powerSource")
	if err != nil {
		context.logger.WithError(err).Error("could not create power source link")
		return nil, ErrInternalServer
	}
	v.AddLink(link)

	link, err = s.manufacturerLink(context, v.Manufacturer, "manufacturer")
	if err != nil {
		context.logger.WithError(err).Error("could not create manufacturer link")
		return nil, ErrInternalServer
	}
	v.AddLink(link)

	return v, nil
}

// GetPowerSources gets all available power sources.
func (s *Service) GetPowerSources(context *Context) (interface{}, error) {

	context.logger.Infof("get power sources")

	entities, err := s.repository.GetPowerSources()
	if err != nil {
		if context.server.IsCriticalError(err) {
			context.logger.WithError(err).Error("could not get power sources")
			return nil, ErrInternalServer
		}
		return nil, err
	}
	for _, m := range entities {
		link, err := s.powerSourceLink(context, m, "canonical")
		if err != nil {
			context.logger.WithError(err).Error("could not create power source link")
			return nil, ErrInternalServer
		}
		m.AddLink(link)
	}

	return entities, nil

}

// GetPowerSource gets the specified power source.
func (s *Service) GetPowerSource(context *Context) (interface{}, error) {

	id := context.Params["id"]

	context.logger.Infof("get power source by id: '%s'", id)

	p, err := s.repository.GetPowerSource(id)
	if err != nil {
		if context.server.IsCriticalError(err) {
			context.logger.WithError(err).Error("could not get power source")
			return nil, ErrInternalServer
		}
		return nil, err
	}

	return p, nil
}
