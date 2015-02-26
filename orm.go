package kwiscale

// ORM interface system. IMHO Related and Join are problematic
// because I really want to have a Mongo Plugin. Mongo fetch the whole
// mapped object... So, orm plugins for other database drivers should
// make manipulation to force them to be set. Example:
//
// type Email struct {
//		For string
//		Url string
// }
//
// type User struct {
//		Name sting
//		Emails []Email // -> mapped to Email
// }
//
// Fetching IORM.First(&u, User{Name:"foo"}) should set u.Emails
// At this time, we have to use:
//
// emails := make([]Email)
// u := User{}
// IORM.First(&u, User{Name: "foo"}).Join(&emails)
//
// That's not trivial...

var ormDriverRegistry = make(map[string]IORM)

type IORM interface {
	// ConnectionString declare the connection url (or path) to database
	ConnectionString(string)

	// Init could initiate tables, migration, and so on
	Init(tables ...interface{}) error

	// First should return the first element
	First(elem interface{}, where ...interface{}) IORM

	// Find should return dataset (list)
	Find(elem interface{}, where ...interface{}) IORM

	// Related should map "related" object from "base",
	// that handles relationship
	// Example:
	//		var emails []Email
	//		IORM.Related(&User{Id: 55}, &emails)
	// That set emails to be a list of email data where
	// email user_id is 55
	Related(base, related interface{}) IORM

	// Join is a common way to use Related...
	// Example:
	//		var emails []Email
	//		var user User{}
	//		IROM.Find(&user, User{Id:55}).Join(&emails)
	// emails should be set for email that have user_id to 55
	Join(to interface{}) IORM

	// Save should create or update data
	Save(elem interface{}) IORM

	Limit(int) IORM
	Offset(int) IORM
	Order(string) IORM

	// Set a debug status (if any)
	Debug(bool)
}

// RegisterORMDriver should be called by ORM plugins
// to register a driver in framework.
func RegisterORMDriver(name string, db IORM) {
	ormDriverRegistry[name] = db
}
