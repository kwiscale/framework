package kwiscale

import (
	"reflect"
	"time"
)

/*
Must work like this:

user := &User{}

kwiscale.Datastore().Get(map[string]interface{}{
	"name" : "Foo"
}).Limit(10).Find(u)

*/

var dbdrivers = make(map[string]reflect.Type)

func RegisterDatabase(name string, ds DB) {
	dbdrivers[name] = reflect.ValueOf(ds).Elem().Type()
}

type DBModelable interface {
	OnCreate()
	OnUpdate()
	OnGetFunc(func(DBModelable, interface{}))
	OnGet(res interface{})
	Id(id ...interface{}) interface{}
}

type DBModel struct {
	ID        interface{}
	CreatedAt time.Time
	UpdatedAt time.Time
	onGet     func(DBModelable, interface{})
}

// Id Get/Set result id in m.ID
func (m *DBModel) Id(id ...interface{}) interface{} {
	if len(id) > 0 {
		if len(id) > 1 {
			panic("You must give only one ID to Id() function")
		}

		m.ID = id[0]
	}
	return m.ID
}

func (model *DBModel) OnGetFunc(f func(DBModelable, interface{})) {
	model.onGet = f
}

func (model *DBModel) OnGet(res interface{}) {
	model.onGet(model, res)
}

func (model *DBModel) OnCreate() {
	model.CreatedAt = time.Now()
}

func (model *DBModel) OnUpdate() {
	model.UpdatedAt = time.Now()
}

type DBOptions map[string]interface{}

type Q map[string]interface{}

type DB interface {
	SetOptions(DBOptions)
	Init()
	Insert(what interface{}) error
	Get(id, result interface{})
	Find(query Q) DBQuery
	Update(where Q, what interface{}) error
	Delete(where Q) error
	Close()
}

type DBQuery interface {
	Limit(int64) DBQuery
	Offset(int64) DBQuery
	One(interface{}) error
	All([]interface{}) error
}
