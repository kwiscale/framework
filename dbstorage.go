package kwiscale

/*Currently disabled*/

/*
import (
	"reflect"
	"time"
)

var dbdrivers = make(map[string]reflect.Type)

func RegisterDatabase(name string, ds DB) {
	dbdrivers[name] = reflect.ValueOf(ds).Elem().Type()
}

// Database options to give to driver.
type DBOptions map[string]interface{}

// Query structure.
type Q struct {
	Request interface{}
	Limit   int
	Offset  int
}

type DB interface {
	SetOptions(DBOptions)
	Init()
	Close()

	// Migrate can be called to migrate table.
	Migrate(interface{}) error
	// Save should update/create data.
	Save(obj interface{}) error
	// Fetch queries database to fetch data and set it on result.
	Fetch(request Q, result interface{}) error
	// Delete data following given interface.
	Delete(where interface{}) error
}

type DBModel struct {
	ID        uint `gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}
*/
