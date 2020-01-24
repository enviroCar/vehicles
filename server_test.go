package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-pg/pg/v9"
)

func BuildTestServer(t *testing.T) (*Server, func() error, func() error) {
	t.Log("init new repository")
	repository := NewRepository(&pg.Options{
		User:     getenv("DB_USER", "postgres"),
		Password: getenv("DB_PASS", "postgres"),
		Database: getenv("DB_NAME", "vehicles"),
		Addr:     getenv("DB_ADDR", "localhost:5432"),
	})

	t.Log("init new service")
	service := NewService(repository)

	t.Log("init new server and add routes")
	server := NewServer()

	server.Get("/", service.GetRoot)
	server.Get("/manufacturers", service.GetManufacturers)
	server.Get("/manufacturers/{hsn}", service.GetManufacturer)
	server.Get("/manufacturers/{hsn}/vehicles", service.GetVehicles)
	server.Get("/manufacturers/{hsn}/vehicles/{tsn}", service.GetVehicle)
	server.Get("/powerSources", service.GetPowerSources)
	server.Get("/powerSources/{id}", service.GetPowerSource)

	return server, repository.Close, service.Close
}

func TestServerRoot(t *testing.T) {

	server, repositoryClose, serviceClose := BuildTestServer(t)
	defer repositoryClose()
	defer serviceClose()

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

	t.Log("get root")
	server.ServeHTTP(rr, req)

	AssertOkStatusCode(t, rr.Code)

	want := `{"links":[{"href":"https://processing.envirocar.org/vehicles/manufacturers","type":"application/json","title":"Manufacturers","rel":"manufacturers"},{"href":"https://processing.envirocar.org/vehicles/powerSources","type":"application/json","title":"Power Sources","rel":"powerSources"}]}`
	AssertResponseBody(t, rr.Body.String(), want)

	t.Logf("response body: %v", rr.Body.String())
}

func TestServerGetManufacturerById(t *testing.T) {

	server, repositoryClose, serviceClose := BuildTestServer(t)
	defer repositoryClose()
	defer serviceClose()

	req, err := http.NewRequest("GET", "/manufacturers/0005", nil)
	if err != nil {
		t.Fatal(err)
	}

	req.Host = "localhost"
	req.Host = "processing.envirocar.org"
	req.Header.Add("Host", "processing.envirocar.org")
	req.Header.Add("accept", "application/json")

	rr := httptest.NewRecorder()

	t.Log("get manufacturer by id")
	server.ServeHTTP(rr, req)

	AssertOkStatusCode(t, rr.Code)

	want := `{"links":[{"href":"http://processing.envirocar.org/manufacturers/0005/vehicles","type":"application/json","rel":"vehicles"},{"href":"http://processing.envirocar.org/manufacturers/0005","type":"application/json","title":"BMW","rel":"self"}],"hsn":"0005","name":"BMW"}`
	AssertResponseBody(t, rr.Body.String(), want)

	t.Logf("response body: %v", rr.Body.String())
}

func TestServerGetPowerSourceById(t *testing.T) {

	server, repositoryClose, serviceClose := BuildTestServer(t)
	defer repositoryClose()
	defer serviceClose()

	req, err := http.NewRequest("GET", "/powerSources/14", nil)
	if err != nil {
		t.Fatal(err)
	}

	req.Host = "localhost"
	req.Host = "processing.envirocar.org"
	req.Header.Add("Host", "processing.envirocar.org")
	req.Header.Add("accept", "application/json")

	rr := httptest.NewRecorder()

	t.Log("get power source by id")
	server.ServeHTTP(rr, req)

	AssertOkStatusCode(t, rr.Code)

	want := `{"id":14,"name":"Wasserst./Benzin/E","description":"Bivalenter Betrieb mit Wasserstoff oder Benzin kombiniert mit Elektromotor"}`
	AssertResponseBody(t, rr.Body.String(), want)

	t.Logf("response body: %v", rr.Body.String())
}

func TestServerGetVehicleByManufacturer(t *testing.T) {

	server, repositoryClose, serviceClose := BuildTestServer(t)
	defer repositoryClose()
	defer serviceClose()

	req, err := http.NewRequest("GET", "/manufacturers/0005/vehicles/155", nil)
	if err != nil {
		t.Fatal(err)
	}

	req.Host = "localhost"
	req.Host = "processing.envirocar.org"
	req.Header.Add("Host", "processing.envirocar.org")
	req.Header.Add("accept", "application/json")

	rr := httptest.NewRecorder()

	t.Log("get vehicle by manufacturer and vehicle id")
	server.ServeHTTP(rr, req)

	AssertOkStatusCode(t, rr.Code)
	want := `{"links":[{"href":"http://processing.envirocar.org/manufacturers/0005/vehicles/155","type":"application/json","title":"645CI","rel":"self"},{"href":"http://processing.envirocar.org/powerSources/1","type":"application/json","title":"Benzin","rel":"powerSource"},{"href":"http://processing.envirocar.org/manufacturers/0005","type":"application/json","title":"BMW","rel":"manufacturer"}],"tsn":"155","commercialName":"645CI","allotmentDate":"2003-07-01","category":"01","bodywork":"0200","power":245,"engineCapacity":4398,"axles":2,"poweredAxles":1,"seats":4,"maximumMass":2070}`
	AssertResponseBody(t, rr.Body.String(), want)

	t.Logf("response body: %v", rr.Body.String())
}

func TestServerGetManufacturers(t *testing.T) {

	server, repositoryClose, serviceClose := BuildTestServer(t)
	defer repositoryClose()
	defer serviceClose()

	req, err := http.NewRequest("GET", "/manufacturers", nil)
	if err != nil {
		t.Fatal(err)
	}

	req.Host = "localhost"
	req.Host = "processing.envirocar.org"
	req.Header.Add("Host", "processing.envirocar.org")
	req.Header.Add("accept", "application/json")

	rr := httptest.NewRecorder()

	t.Log("get manufacturers")
	server.ServeHTTP(rr, req)

	AssertOkStatusCode(t, rr.Code)
}

func TestServerGetVehiclesByManufacturer(t *testing.T) {

	server, repositoryClose, serviceClose := BuildTestServer(t)
	defer repositoryClose()
	defer serviceClose()

	req, err := http.NewRequest("GET", "/manufacturers/0005/vehicles", nil)
	if err != nil {
		t.Fatal(err)
	}

	req.Host = "localhost"
	req.Host = "processing.envirocar.org"
	req.Header.Add("Host", "processing.envirocar.org")
	req.Header.Add("accept", "application/json")

	rr := httptest.NewRecorder()

	t.Log("get vehicles by manufacturer")
	server.ServeHTTP(rr, req)

	AssertOkStatusCode(t, rr.Code)
}

func TestServerGetPowerSources(t *testing.T) {

	server, repositoryClose, serviceClose := BuildTestServer(t)
	defer repositoryClose()
	defer serviceClose()

	req, err := http.NewRequest("GET", "/powerSources", nil)
	if err != nil {
		t.Fatal(err)
	}

	req.Host = "localhost"
	req.Host = "processing.envirocar.org"
	req.Header.Add("Host", "processing.envirocar.org")
	req.Header.Add("accept", "application/json")

	rr := httptest.NewRecorder()

	t.Log("get power sources")
	server.ServeHTTP(rr, req)

	AssertOkStatusCode(t, rr.Code)
}

func AssertOkStatusCode(t *testing.T, code int) {
	// Check the status code is what we expect.
	if code != http.StatusOK {
		t.Fatalf("handler returned wrong status code: got %v want %v",
			code, http.StatusOK)
	}
}

func AssertResponseBody(t *testing.T, got, want string) {
	// Check the response body is what we expect.
	if strings.TrimSpace(got) != want {
		t.Fatalf("handler returned unexpected body: got %v want %v",
			got, want)
	}
}
