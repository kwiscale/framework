package db

// initialize dbs drivers
import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

// InitDB returns usable sql.DB for giver driver name and db uri
func InitDB(driver, uri string) (*gorm.DB, error) {
	db, err := gorm.Open(driver, uri)
	if err != nil {
		return nil, err
	}

	// Get database connection handle [*gorm.DB](http://golang.org/pkg/database/sql/#DB)
	return &db, nil
}
