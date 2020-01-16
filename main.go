package main

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"

	"github.com/go-pg/pg/v9"
	_ "github.com/go-pg/pg/v9/orm"
)

func main() {

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

	re, err := regexp.Compile("^(?P<v0>[^.]+)$")
	if err != nil {
		log.Fatal(err)
	}
	log.Println(re.MatchString("localhost"))

	log.Fatal(server.Start(fmt.Sprintf(":%d", getPort())))
}

func getenv(name, defaultValue string) string {
	if value, ok := os.LookupEnv(name); ok {
		return value
	}
	return defaultValue
}

func getPort() int {
	portStr := os.Getenv("PORT")
	if portStr != "" {
		if port, err := strconv.Atoi(portStr); err == nil {
			return port
		}
	}
	return 8080
}
