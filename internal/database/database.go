package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
	"sync"
)

// Helpers
const (
	InitSQLScriptPath string = "/internal/database/sqlscripts/db_init.sql"
	INITsql           string = `
		CREATE TABLE IF NOT EXISTS roles 
		(
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT UNIQUE
		);
		
CREATE TABLE IF NOT EXISTS member_roles 
(
	memberid INTEGER,
	roleid INTEGER,
	PRIMARY KEY (memberid, roleid),
	FOREIGN KEY (memberid) REFERENCES members (id),
	FOREIGN KEY (roleid) REFERENCES roles (id)
);
CREATE TABLE IF NOT EXISTS members 
(
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	guild TEXT,
	username TEXT,
	nickname TEXT,
	avatarurl TEXT,
	displayavatarurl TEXT,
	bannerurl TEXT,
	displaybannerurl TEXT,
	usercolor TEXT,
	joinedat TEXT,
	userstatus TEXT,
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

CREATE TABLE IF NOT EXISTS lines (
	lineid integer primary key AUTOINCREMENT,
	bid integer,
	phrase text,
	author varchar(255),
	toid varchar(255),
	ltype varchar(255),
	createdat DATETIME DEFAULT CURRENT_TIMESTAMP,
	foreign key (bid) references bots (botid)
);
`
)

var (
	AllowedMemberCols = map[string]bool{
		// Members
		"ids":        true,
		"guilds":     true,
		"usernames":  true,
		"nicknames":  true,
		"avatarurls": true,
		"usercolors": true,
		"msgcounts":  true,
		"joinedats":  true,
	}

	AllowedBotCols = map[string]bool{
		// Bots
		"botids":   true,
		"botnames": true,
		"authors":  true,
	}

	AllowedLineCols = map[string]bool{
		// Lines
		"lineids": true,
		"toids":   true,
		"bids":    true,
		"phrases": true,
		"ltypes":  true,
	}
)

// sqlite path db file
func InitDB(dbpath string, initScript string) error {
	// Create and init the database!
	var err error
	var logBuilder strings.Builder
	logBuilder.WriteString("Initializing Database...\n")
	logBuilder.WriteString(fmt.Sprintf("[INIT DB] Path: %s\n", dbpath))

	dbh := GetConnector(dbpath)
	err = dbh.openConnection()
	if err != nil {
		log.Print("Error initializing database connection..., will continue in mem: " + err.Error())
		return err
	}

	defer dbh.DB.Close()

	fileContent, err := os.ReadFile(initScript)
	if err != nil {
		logBuilder.WriteString("Could not read from file, will use a default instead...: \n")
		logBuilder.WriteString("[INIT DB] Running default script...\n")
		res, err := dbh.RunSQLscript(INITsql)
		if err != nil {
			log.Print("[INIT DB]Error initializing the database: " + err.Error())
		}
		logBuilder.WriteString(fmt.Sprintf("SQL script result: %s\n", res))

	} else {
		logBuilder.WriteString("[INIT DB] Running script from file...\n")
		res, err := dbh.RunSQLscript(string(fileContent))
		if err != nil {
			log.Print("[INIT DB]Error initializing the database: " + err.Error())
		}
		logBuilder.WriteString(fmt.Sprintf("[INIT DB] Script execution result: %s\n", res))

	}
	logBuilder.WriteString("\n")
	log.Print(logBuilder.String())

	return err
}

type DBHandler struct {
	DB     *sql.DB
	dbFile string
	mu     sync.Mutex
}

// Opens a connection to the database and holds reference to the Struct Handler
// Should Close The Connection!
func GetConnector(dbPath string) *DBHandler {
	return &DBHandler{dbFile: dbPath}
}

func (dbh *DBHandler) openConnection() error {
	var err error

	dbh.DB, err = sql.Open("sqlite", dbh.dbFile)
	if err != nil {
		return err
	}

	return nil
}

func (dbh *DBHandler) Close() {
	if dbh.DB != nil {
		dbh.DB.Close()
	}
}

func (dbh *DBHandler) Metrics(mtype string) (string, error) {
	mu := &dbh.mu
	mu.Lock()
	defer mu.Unlock()

	err := dbh.openConnection()
	if err != nil {
		log.Printf("There's been an error creating the DB handler: %v", err)
		return "Failed to create DB handler", err
	}
	var result string

	switch mtype {
	case "all":
		// Decide what metrics to get.
	case "members":

	case "bots":

	case "lines":

	default:
		log.Print("Shouldn't be here...")
		return "", fmt.Errorf("Incorrect type ")
	}

	return result, nil
}

// Sql execution funcs
// Should be used to initialize the database table
func (dbh *DBHandler) RunSQLscript(sql string) (string, error) {
	dbh.mu.Lock()
	defer dbh.mu.Unlock()

	err := dbh.openConnection()
	if err != nil {
		log.Printf("There's been an error opening the DB connection." + err.Error())
		return "No result...", err
	}
	defer dbh.Close()

	result, err := dbh.DB.Exec(sql)
	if err != nil {
		return "No resutt...", fmt.Errorf("error executing SQL script: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Sprintf("Script executed. Rows affected error: " + err.Error()), err
	}

	return fmt.Sprintf("Script executed. Rows affected: %v", rowsAffected), nil
}

// Health checking funcs
type SchemaInfo struct {
	Tables []TableInfo `json:"tables"`
}

// TableInfo holds information about a single table.
type TableInfo struct {
	Name    string       `json:"name"`
	Columns []ColumnInfo `json:"columns"`
	Indexes []IndexInfo  `json:"indexes"`
}

// ColumnInfo holds information about a single column in a table.
type ColumnInfo struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

// IndexInfo holds information about an index in a table.
type IndexInfo struct {
	Name  string `json:"name"`
	Index string `json:"index"`
}

func DBHealthCheck(dbpath string) *SchemaInfo {
	// Open the database connection
	dbh := GetConnector(dbpath)

	dbh.mu.Lock()
	defer dbh.mu.Unlock()

	err := dbh.openConnection()
	if err != nil {
		log.Print("Error opening connection to the Database")
		return nil
	}
	defer dbh.Close()

	// Get the list of tables
	rows, err := dbh.DB.Query(`SELECT name FROM sqlite_master WHERE type='table'`)
	if err != nil {
		log.Print("Failed to query tables")
		return nil
	}
	defer rows.Close()

	var schema SchemaInfo

	// Iterate through the tables
	for rows.Next() {
		var tableName string
		if err := rows.Scan(&tableName); err != nil {
			log.Print("Failed to scan table name")
			return nil
		}

		// log.Printf("Table name: %s\n", tableName)

		// Get the columns for each table
		columnsRows, err := dbh.DB.Query(`PRAGMA table_info(` + tableName + `)`)
		if err != nil {
			log.Print("Failed to get table info.")
			return nil
		}
		defer columnsRows.Close()

		var columns []ColumnInfo
		for columnsRows.Next() {
			var column ColumnInfo
			var cid int
			var notNull, pk int
			var dfltValue interface{}
			if err := columnsRows.Scan(&cid, &column.Name, &column.Type, &notNull, &dfltValue, &pk); err != nil {
				log.Print("Failed to scan column info")
				return nil
			}
			// log.Printf("Column info %+v\n", column)
			columns = append(columns, column)
		}

		// Get the indexes for the table
		indexRows, err := dbh.DB.Query(`PRAGMA index_list(` + tableName + `)`)
		if err != nil {
			log.Print("Failed to get indexes")
			return nil
		}
		defer indexRows.Close()

		var indexes []IndexInfo

		for indexRows.Next() {
			var index IndexInfo
			var seq int // Sequence number
			var unique int
			var origin string
			var partial int

			// Correct order of scanning the columns from PRAGMA index_list()
			if err := indexRows.Scan(&seq, &index.Name, &unique, &origin, &partial); err != nil {
				log.Printf("Failed to scan index info for table %s: %v", tableName, err)
				return nil
			}

			// Get index info
			indexDetailsRows, err := dbh.DB.Query(`PRAGMA index_info(` + index.Name + `)`)
			if err != nil {
				log.Printf("Failed to get index details for index %s: %v", index.Name, err)
				return nil
			}
			defer indexDetailsRows.Close()

			var colNames []string
			for indexDetailsRows.Next() {
				var seqNo, cid int
				var colName string
				if err := indexDetailsRows.Scan(&seqNo, &cid, &colName); err != nil {
					log.Print("Failed to scan index details")
					return nil
				}
				colNames = append(colNames, colName)
			}
			index.Index = fmt.Sprintf("%v", colNames)
			indexes = append(indexes, index)
		}

		// Add the table info to the schema
		schema.Tables = append(schema.Tables, TableInfo{
			Name:    tableName,
			Columns: columns,
			Indexes: indexes,
		})
	}

	// Return the schema info
	return &schema
}

// utils
func isNumeric(s string) bool {
	re := regexp.MustCompile(`^[0-9]+$`)
	return re.MatchString(s)
}

func interfaceSlice(slice []string) []interface{} {
	interfaces := make([]interface{}, len(slice))
	for i, v := range slice {
		interfaces[i] = v
	}
	return interfaces
}
