package main

import (
	"fmt"
	"github.com/go-pg/pg/v9"
	"testing"
)

func NewTesRepository(t *testing.T) *Repository{
	t.Log("create test repository")
	repository := NewRepository(&pg.Options{
		User:     getenv("DB_USER", "postgres"),
		Password: getenv("DB_PASS", "postgres"),
		Database: getenv("DB_NAME", "vehicles"),
		Addr:     getenv("DB_ADDR", "localhost:5432"),
	})
	return repository
}

func TestGetVehicle(t *testing.T) {

	r := NewTesRepository(t)

	m := &Manufacturer{
		ID:     "0005",
		Name:   "BMW",
	}

	t.Log("get vehicle by id and manufacturer")
	v, err := r.GetVehicle(m, "156")
	if err != nil {
		t.Error(err)
	}

	if v == nil {
		t.Fatal("vehicle is nil")
	}
}


func TestGetVehicleNotFound(t *testing.T) {

	r := NewTesRepository(t)

	m := &Manufacturer{
		ID:     "000x",
		Name:   "BMW",
	}

	t.Log("get vehicle by id and manufacturer")
	_, err := r.GetVehicle(m, "156")
	if err != ErrNotFound {
		t.Fatal(fmt.Sprintf("%v, %T", err, err))
	}
}