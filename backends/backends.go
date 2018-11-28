package backends

import (
	"livego/backends/models"
	"livego/backends/yamlbackend"
	"livego/config"
	"log"
)

// Backend :This is the interface that a backend must implement
type Backend interface {
	Init()
	SelectAll() map[string]models.Route
	Select(id string) (models.Route, bool)
	Insert(r models.Route) (models.Route, error)
	Update(id string, r models.Route) (models.Route, error)
	Delete(id string) error
}

var DB Backend

// SetBackend :This function sets the backend db to be used with the program
func SetBackend() {

	log.Println("Database Loaded")
	switch config.Data.Backend["name"] {
	case "yaml":
		log.Printf("Setting Backend to %s", config.Data.Backend["name"])
		DB = &yamlbackend.DB{}
		DB.Init()
	}
}
