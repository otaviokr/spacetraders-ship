package component

import (
	"io"
	"os"

	"gopkg.in/yaml.v3"
)

// Route is the representation of the route file.
type Route struct {
	Route []RouteStop `yaml:"route"`
	Error Error       `yaml:"error"`
}

// RouteStop is the representation of each stop, its location, what to buy and what to sell.
type RouteStop struct {
	Station string         `yaml:"station"`
	Buy     map[string]int `yaml:"buy"`
	Sell    map[string]int `yaml:"sell"`
}

// ReadRouteFile will read the YAML file with the route definition.
func ReadRouteFile(path string) (*Route, error) {
	// read yaml config (route)
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return ReadRouteDescription(f)
}

// ReadRouteDescription will generate the component.Route instance from the data read from YAML file.
func ReadRouteDescription(data io.Reader) (*Route, error) {
	var routes Route
	decoder := yaml.NewDecoder(data)
	if err := decoder.Decode(&routes); err != nil {
		return nil, err
	}
	return &routes, nil
}
