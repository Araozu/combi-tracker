package main

import (
	"fmt"
	"log"
	"strconv"
	"sync"

	"github.com/labstack/echo"
)

type Coordinates struct {
	X float64
	Y float64
}

// global state
var (
	coords    []Coordinates
	coordsMux sync.Mutex
)

func main() {
	e := echo.New()

	// Handler for receiving coordinates
	e.GET("/track", func(c echo.Context) error {
		// Get query params
		xStr := c.QueryParam("x")
		yStr := c.QueryParam("y")

		x, err := strconv.ParseFloat(xStr, 64)
		if err != nil {
			return c.String(400, "Invalid x")
		}

		y, err := strconv.ParseFloat(yStr, 64)
		if err != nil {
			return c.String(400, "Invalid y")
		}

		// Update coordinates
		coordsMux.Lock()
		coords = append(coords, Coordinates{X: x, Y: y})
		coordsMux.Unlock()

		log.Printf("Received coordinates: %f, %f\n", x, y)

		return c.String(200, "Received coordinates")
	})

	// Handler for getting the last coordinate
	e.GET("/", func(c echo.Context) error {
		coordsMux.Lock()
		defer coordsMux.Unlock()

		coords_len := len(coords)
		if coords_len == 0 {
			return c.NoContent(204)
		}

		coord := coords[coords_len-1]

		return c.String(200, fmt.Sprintf("X: %f, Y: %f", coord.X, coord.Y))
	})

	log.Printf("Hello, world (Starting server)\n")
	e.Logger.Fatal(e.Start(":8888"))
}
