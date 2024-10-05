package database

import (
	"database/sql"
	"log"
	"os"
	"regexp"
	"sync"
)

const (
	InitSQLScriptPath string = "C:\\Users\\kyria\\Documents\\Coding\\Discord Bots\\internal\\database\\db_init.sql"
	INITsql           string = `
		CREATE TABLE IF NOT EXISTS roles 
		(
			id INTEGER PRIMARY KEY AUTOINCREMENT,
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
			id INTEGER PRIMARY KEY AUTOINCREMENT,
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
		);
		-- DROP TABLE bots;
		CREATE TABLE IF NOT EXISTS bots (
			botid integer primary key AUTOINCREMENT,
			botguild varchar(255),
		    botname varchar(255),
			avatarurl varchar(255),
			bannerurl varchar(255),
		    createdat varchar(255),
			author varchar(255),
		    botstatus varchar(255),
		    isSinger boolean
		
		);

		CREATE TABLE IF NOT EXISTS trigger_words (
		    trigger_id integer primary key AUTOINCREMENT,
		    bot_id integer foreign key REFERENCES bots (botid),
		    phrase varchar(255),
		    author char(20),
		    data_datetime DATETIME DEFAULT CURRENT_TIMESTAMP
		);

		CREATE TABLE IF NOT EXISTS typical_lines (
		    line_id integer primary key AUTOINCREMENT,
		    bot_id integer foreign key REFERENCES bots (botid) ,
		    phrase varchar(255),
		    author char(20),
		    data_datetime DATETIME DEFAULT CURRENT_TIMESTAMP
		);

		CREATE TABLE IF NOT EXISTS sad_lines (
		    line_id integer primary key AUTOINCREMENT,
		    bot_id integer foreign key references bots (botid),
		    phrase varchar(255),
		    author char(20),
		    data_datetime DATETIME DEFAULT CURRENT_TIMESTAMP
		);`
)

func isNumeric(s string) bool {
	re := regexp.MustCompile(`^[0-9]+$`)
	return re.MatchString(s)
}

// sqlite path db file
func InitDB(dbpath string, initScript string) error {
	// Create and init the database!
	dbh := GetConnection(dbpath)
	err := dbh.openConnection()
	if err != nil {
		log.Print("Error initializing database connection..., will continue in mem: " + err.Error())
		return err
	}

	defer dbh.DB.Close()

	fileContent, err := os.ReadFile(initScript)
	if err != nil {
		log.Printf("There was an error reading the sql file, will use a default instead...: " + err.Error())
		dbh.RunSQLscript(INITsql)

	} else {
		dbh.RunSQLscript(string(fileContent))
	}

	fileContent, err = os.ReadFile("C:\\Users\\kyria\\Documents\\Coding\\Discord Bots\\internal\\database\\populate_current_tables.sql")
	if err != nil {
		log.Printf("There was an error executing populative script")
	} else {
		dbh.RunSQLscript(string(fileContent))
	}
	return nil

}

type DBHandler struct {
	DB *sql.DB
	MU sync.Mutex

	dbfile string
}

// Opens a connection to the database and holds reference to the Struct Handler
// Should Close The Connection!
func GetConnection(dbpath string) *DBHandler {
	dbh := &DBHandler{}

	dbh.dbfile = dbpath

	return dbh
}

func (dbh *DBHandler) openConnection() error {
	var err error

	dbh.DB, err = sql.Open("sqlite", dbh.dbfile)
	if err != nil {
		return err
	}

	return nil

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
