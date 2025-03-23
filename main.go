package main

import (
	"fmt"
	"log"
	"strconv"
	"sync"

	"github.com/labstack/echo/v4"
)

type Coordinates struct {
	Lat  float64
	Long float64
}

type CoordinatesStore struct {
	coords []Coordinates
	mu     sync.Mutex
}

func (cs *CoordinatesStore) Add(lat, long float64) {
	cs.mu.Lock()
	defer cs.mu.Unlock()
	cs.coords = append(cs.coords, Coordinates{Lat: lat, Long: long})
}

func (cs *CoordinatesStore) Latest() (*Coordinates, bool) {
	cs.mu.Lock()
	defer cs.mu.Unlock()

	clen := len(cs.coords)
	if clen == 0 {
		return nil, false
	}

	return &cs.coords[clen-1], true
}

var store = CoordinatesStore{}

func main() {
	e := echo.New()

	// Handler for receiving coordinates
	e.GET("/track", func(c echo.Context) error {
		// Get query params
		latStr := c.QueryParam("lat")
		longStr := c.QueryParam("long")

		lat, err := strconv.ParseFloat(latStr, 64)
		if err != nil {
			return c.String(400, "Invalid lat")
		}

		long, err := strconv.ParseFloat(longStr, 64)
		if err != nil {
			return c.String(400, "Invalid long")
		}

		// Update coordinates
		store.Add(lat, long)

		log.Printf("Received coordinates: %f, %f\n", lat, long)

		return c.String(200, fmt.Sprintf("Received: Lat %f, Long %f", lat, long))
	})

	// Handler for getting the last coordinate
	e.GET("/", func(c echo.Context) error {
		coord, ok := store.Latest()
		if !ok {
			return c.NoContent(204)
		}

		return c.String(200, fmt.Sprintf("Lat: %f, Long: %f", coord.Lat, coord.Long))
	})

	e.Logger.Fatal(e.Start(":8888"))
}
