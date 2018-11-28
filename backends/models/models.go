package models

// Endpoint : This defines the destination properties of the route
type Endpoint struct {
	ID      string `yaml:"ID"`
	Host    string `yaml:"Host"`
	Port    string `yaml:"Port,omitempty"`
	App     string `yaml:"App,omitempty"`
	Stream  string `yaml:"Stream,omitempty"`
	Enabled bool   `yaml:"Enabled"`
}

// Route : This struct defines a single route
type Route struct {
	Stream    string     `yaml:"Stream,omitempty"`
	CopyKey   bool       `yaml:"CopyKey"`
	Enabled   bool       `yaml:"Enabled"`
	Endpoints []Endpoint `yaml:"Endpoints"`
}
