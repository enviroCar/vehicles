package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-pg/pg/v9"
)

func TestServer(t *testing.T) {

	repository := NewRepository(&pg.Options{
		User:     getenv("DB_USER", "postgres"),
		Password: getenv("DB_PASS", "postgres"),
		Database: getenv("DB_NAME", "vehicles"),
		Addr:     getenv("DB_ADDR", "localhost:5432"),
	})
	defer repository.Close()

	s := NewService(repository)
	defer s.Close()

	server := NewServer()

	server.Get("/", s.GetRoot)
	server.Get("/manufacturers", s.GetManufacturers)
	server.Get("/manufacturers/{hsn}", s.GetManufacturer)
	server.Get("/manufacturers/{hsn}/vehicles", s.GetVehicles)
	server.Get("/manufacturers/{hsn}/vehicles/{tsn}", s.GetVehicle)
	server.Get("/powerSources", s.GetPowerSources)
	server.Get("/powerSources/{id}", s.GetPowerSource)

	// Create a request to pass to our handler. We don't have any query parameters for now, so we'll
	// pass 'nil' as the third parameter.
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	req.Host = "localhost"
	req.Host = "processing.envirocar.org"
	req.Header.Add("Host", "processing.envirocar.org")
	req.Header.Add("X-Forwarded-Proto", "https")
	req.Header.Add("X-Forwarded-Prefix", "/vehicles")
	req.Header.Add("X-Forwarded-Port", "443")
	req.Header.Add("accept", "application/json")

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()

	server.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	expected := `{"links":[{"href":"https://processing.envirocar.org/vehicles/manufacturers","type":"application/json","title":"Manufacturers","rel":"manufacturers"},{"href":"https://processing.envirocar.org/vehicles/powerSources","type":"application/json","title":"Power Sources","rel":"powerSources"}]}`
	if strings.TrimSpace(rr.Body.String()) != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}
