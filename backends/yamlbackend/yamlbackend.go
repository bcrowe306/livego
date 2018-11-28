package yamlbackend

import (
	"io/ioutil"
	"livego/backends/models"
	"log"
	"os"
	"sync"

	"github.com/segmentio/ksuid"
	"gopkg.in/yaml.v2"
)

var once sync.Once

type routes map[string]models.Route

// DB :This struct defines implemets the Backends interface with additional functionality.
type DB struct {
	Routes   map[string]models.Route
	FileName string
}

// Init :This function Initialized the YAML backend
func (y *DB) Init() {
	if y.FileName == "" {
		log.Println("No route file configured. Setting routes file to ./routes.yaml")
		y.FileName = "./routes.yaml"
	}
	y.NewRoutesInstance()
	y.Load()

}

// NewRoutesInstance :Creates a new instance of routes
func (y *DB) NewRoutesInstance() map[string]models.Route {
	once.Do(func() {
		y.Routes = make(routes)
	})
	return y.Routes
}

// Load :Routes :This Function loads routes from the routes file and populates the exposed Routes object
func (y *DB) Load() error {
	log.Printf("Loading routes from %s", y.FileName)
	fileBytes, err := ioutil.ReadFile(y.FileName)
	if err != nil {
		if os.IsNotExist(err) {
			ioutil.WriteFile(y.FileName, []byte{}, 0777)
			return nil
		}
		return err
	}
	log.Printf("Parsing routes from %s", y.FileName)
	yaml.Unmarshal(fileBytes, y.Routes)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (y *DB) save() error {
	rts := y.NewRoutesInstance()
	outBytes, err := yaml.Marshal(rts)
	if err != nil {
		log.Println(err)
		return err
	}
	ioutil.WriteFile(y.FileName, outBytes, 0777)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

// Insert :This function inserts a route unto the Routes object and persist it to storage
func (y DB) Insert(r models.Route) (models.Route, error) {
	rts := y.NewRoutesInstance()
	for i := range r.Endpoints {
		r.Endpoints[i].ID = ksuid.New().String()
	}
	routeID := ksuid.New().String()
	rts[routeID] = r
	err := y.save()
	if err != nil {
		return models.Route{}, err
	}
	return rts[routeID], nil
}

// SelectAll :This function selects  all routes
func (y DB) SelectAll() map[string]models.Route {
	log.Println("Select All Routes")
	return y.Routes
}

// Select :This function selects a route based on App, and potentially any other matching value
func (y DB) Select(id string) (models.Route, bool) {
	rts := y.NewRoutesInstance()
	if r, exist := rts[id]; exist {
		return r, true
	}
	return models.Route{}, false
}

// Delete :This function deletes the element from the Routes array by the ID listed in the arguments.
func (y DB) Delete(id string) error {
	rts := y.NewRoutesInstance()
	if _, exist := rts[id]; exist {
		delete(y.Routes, id)
	}
	if err := y.save(); err != nil {
		return err
	}
	return nil
}

// Update :This function upserts the Route in the Routes map by ID, and save it to the file specified
func (y DB) Update(id string, newRoute models.Route) (models.Route, error) {
	rts := y.NewRoutesInstance()
	if _, exist := rts[id]; exist {
		rts[id] = newRoute
		if err := y.save(); err != nil {
			return models.Route{}, err
		}
		return rts[id], nil
	}
	newID := ksuid.New().String()
	rts[newID] = newRoute
	if err := y.save(); err != nil {
		return models.Route{}, err
	}
	return rts[newID], nil
}
