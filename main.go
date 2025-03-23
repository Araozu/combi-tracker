package main

import (
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/labstack/echo/v4"
)

type Coordinates struct {
	Lat        float64
	Long       float64
	ClientTime int64
	ServerTime int64
}

type CoordinatesStore struct {
	coords []Coordinates
	mu     sync.Mutex
}

func (cs *CoordinatesStore) Add(lat, long float64, clientTime int64) *Coordinates {
	cs.mu.Lock()
	defer cs.mu.Unlock()

	newCoor := Coordinates{
		Lat:        lat,
		Long:       long,
		ClientTime: clientTime,
		ServerTime: time.Now().UnixMilli(),
	}

	cs.coords = append(cs.coords, newCoor)
	return &newCoor
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
		timeStr := c.QueryParam("time")

		lat, err := strconv.ParseFloat(latStr, 64)
		if err != nil {
			return c.String(400, "Invalid lat")
		}

		long, err := strconv.ParseFloat(longStr, 64)
		if err != nil {
			return c.String(400, "Invalid long")
		}

		time, err := strconv.ParseInt(timeStr, 10, 64)
		if err != nil {
			return c.String(400, "Invalid unix time (ms)")
		}

		// Update coordinates
		inserted := store.Add(lat, long, time)

		return c.String(200, fmt.Sprintf(
			"Lat %v, Long %v, Client time %d, Server time %d, Time delta %d ms",
			inserted.Lat,
			inserted.Long,
			inserted.ClientTime,
			inserted.ServerTime,
			inserted.ServerTime-inserted.ClientTime,
		))
	})

	// Handler for getting the last coordinate
	e.GET("/", func(c echo.Context) error {
		coord, ok := store.Latest()
		if !ok {
			return c.NoContent(204)
		}

		return c.String(200, fmt.Sprintf("Lat: %f, Long: %f", coord.Lat, coord.Long))
	})

	e.GET("/timestamp", func(c echo.Context) error {
		return c.String(200, fmt.Sprintf("%d", time.Now().UnixMilli()))
	})

	e.Logger.Fatal(e.Start(":8888"))
}
