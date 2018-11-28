package backends

import (
	"livego/backends/models"
	"livego/backends/yamlbackend"
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
func SetBackend(bkend string) {
	log.Printf("Setting Backend to %s", bkend)
	DB = &yamlbackend.DB{}
	DB.Init()

	log.Println("Database Loaded")
}
