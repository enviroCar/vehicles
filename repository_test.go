package main

import (
	"fmt"
	"github.com/go-pg/pg/v9"
	"net/http"
	"testing"
)

func NewTestRepository(t *testing.T) *Repository {
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

	r := NewTestRepository(t)

	m := &Manufacturer{
		ID:   "0005",
		Name: "BMW",
	}

	t.Log("get vehicle by id and manufacturer")
	v, err := r.GetVehicle(m, "156")
	if err != nil {
		t.Fatal(err)
	}

	if v == nil {
		t.Fatal("vehicle is nil")
	}
	t.Log(v)
}

func TestGetVehicleNotFound(t *testing.T) {

	r := NewTestRepository(t)

	m := &Manufacturer{
		ID:   "000x",
		Name: "BMW",
	}

	t.Log("get vehicle by id and manufacturer")
	_, err := r.GetVehicle(m, "156")
	if err != ErrNotFound {
		t.Fatal(fmt.Sprintf("%v, %T", err, err))
	}
}

func TestGetManufacturer(t *testing.T) {

	r := NewTestRepository(t)

	t.Log("get manufacturer")
	m, err := r.GetManufacturer("0005")
	if err != nil {
		t.Fatal(err)
	}

	if m == nil {
		t.Fatal("manufacturer is nil")
	}
	t.Log(m)
}

func TestGetManufacturerNotFound(t *testing.T) {

	r := NewTestRepository(t)

	t.Log("get manufacturer")
	_, err := r.GetManufacturer("000x")
	if err != ErrNotFound {
		t.Fatal(fmt.Sprintf("%v, %T", err, err))
	}
}

func TestGetPowerSourceParserError(t *testing.T) {

	r := NewTestRepository(t)

	t.Log("get power source by id")
	_, err := r.GetPowerSource("1x")
	if err == nil {
		t.Fatal("error is nil")
	}

	httpError, ok := err.(Error)
	if !ok {
		t.Fatal(fmt.Sprintf("%v, %T", err, err))
	}

	if httpError.Status() != http.StatusNotFound {
		t.Fatalf("status code is bad, got:'%v', want:'%v'", httpError.Status(), http.StatusNotFound)
	}

	t.Log(fmt.Sprintf("%v, %T", err, err))
}
