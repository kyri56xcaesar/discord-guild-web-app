package servicedb

import (
	"database/sql"
	"log"
	"sync"

	_ "modernc.org/sqlite"
)

type DBHandler struct {
	DB *sql.DB
	MU sync.Mutex
}

func (dbH *DBHandler) OpenConnection(dbname string) (*DBHandler, error) {
	db, err := sql.Open("sqlite", dbname)
	if err != nil {
		return nil, err
	}

	dbH.DB = db

	return dbH, nil
}

// Should be used to initialize the database table
func (dbH *DBHandler) RunSQLscript(sql string) error {

	if dbH.DB == nil {
		log.Print("Must initialize the DB connector")
	}

	db := dbH.DB
	defer db.Close()

	dbH.MU.Lock()
	defer dbH.MU.Unlock()

	log.Println("Running sql script...")
	statement, err := db.Prepare(sql)
	if err != nil {
		return err
	}

	statement.Exec()
	log.Println("sql script executed.")

	return nil
}

const INITsql string = `
		CREATE TABLE IF NOT EXISTS roles 
		(
			id INTEGER PRIMARY KEY,
			name TEXT UNIQUE
		);
		CREATE TABLE IF NOT EXISTS user_roles 
		(
			user_id INTEGER,
			role_id INTEGER,
			PRIMARY KEY (user_id, role_id),
			FOREIGN KEY (user_id) REFERENCES users (id),
			FOREIGN KEY (role_id) REFERENCES roles (id)
		);
		CREATE TABLE IF NOT EXISTS users 
		(
			id INTEGER PRIMARY KEY,
			username TEXT,
			guild TEXT,
			nickname TEXT,
			avatarurl TEXT,
			displayavatar TEXT,
			bannerurl TEXT,
			displaybanner TEXT,
			usercolor TEXT,
			joinedat TEXT,
			status TEXT,
			msgcount INTEGER
		);`
