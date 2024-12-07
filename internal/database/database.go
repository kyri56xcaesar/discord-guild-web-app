package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"

	"kyri56xcaesar/discord-guild-web-app/internal/utils"
)

// Helpers
const (
	TypeMember        string = "members"
	TypeBot           string = "bots"
	TypeLine          string = "lines"
	TypeRole          string = "roles"
	TypeMemberRoles   string = "member_roles"
	InitSQLScriptPath string = "/internal/database/sqlscripts/db_init.sql"
	INITsql           string = `
CREATE TABLE IF NOT EXISTS roles 
(
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	userid INTEGER,
	rolename TEXT,
	rolecolor TEXT,
	foreign key (userid) references members (id) 
);
		
CREATE TABLE IF NOT EXISTS members 
(
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	guild TEXT,
	username TEXT,
	nickname TEXT,
  leaguename TEXT,
	avatarurl TEXT,
	displayavatarurl TEXT,
	bannerurl TEXT,
	displaybannerurl TEXT,
	usercolor TEXT,
	joinedat TEXT,
	status TEXT,
	msgcount INTEGER
);

-- DROP TABLE bots;
CREATE TABLE IF NOT EXISTS bots (
	id integer primary key AUTOINCREMENT,
	guild varchar(255),
  username varchar(255),
	avatarurl varchar(255),
	bannerurl varchar(255),
  createdat varchar(255),
	author varchar(255),
  status varchar(255),
  issinger boolean

);

CREATE TABLE IF NOT EXISTS lines (
	id integer primary key AUTOINCREMENT,
	bid integer,
	phrase text,
	author varchar(255),
	toid varchar(255),
	ltype varchar(255),
	createdat DATETIME DEFAULT CURRENT_TIMESTAMP,
	foreign key (bid) references bots (botid)
);

CREATE TABLE IF NOT EXISTS messages (
	messageid	integer primary key AUTOINCREMENT,
	userid		integer,
	content		text,
	channel		text,
	createdat	text,
	foreign key (userid) references members (id)
);
`
)

type DBHandler struct {
	DB     *sql.DB
	dbFile string
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

var (
	DefaultLimit int    = 100
	DefaultOrder string = "desc"

	AllKeys []string

	AllowedMemberCols = map[string]bool{
		// Members
		"id":         true,
		"guild":      true,
		"username":   true,
		"nickname":   true,
		"leaguename": true,
		"avatarurl":  true,
		"usercolor":  true,
		"msgcount":   true,
		"joinedat":   true,
		"status":     true,
	}

	AllowedBotCols = map[string]bool{
		// Bots
		"id":       true,
		"username": true,
		"author":   true,
	}

	AllowedLineCols = map[string]bool{
		// Lines
		"id":     true,
		"toid":   true,
		"bid":    true,
		"phrase": true,
		"ltype":  true,
	}
)

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

	AllKeys = utils.AppendKeys([]map[string]bool{AllowedMemberCols, AllowedBotCols, AllowedLineCols})
	logBuilder.WriteString("[INIT DB] Key fields initialized.\n")

	log.Print(logBuilder.String())
	return err
}

// Sql execution funcs
// Should be used to initialize the database table
func (dbh *DBHandler) RunSQLscript(sql string) (string, error) {
	err := dbh.openConnection()
	if err != nil {
		log.Printf("There's been an error opening the DB connection." + err.Error())
		return "No result...", err
	}
	defer dbh.Close()

	result, err := dbh.DB.Exec(sql)
	if err != nil {
		return "No result...", fmt.Errorf("error executing SQL script: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Sprintf("Script executed. Rows affected error: " + err.Error()), err
	}

	return fmt.Sprintf("Script executed. Rows affected: %v", rowsAffected), nil
}

// Database check
func DBHealthCheck(dbpath string) *SchemaInfo {
	// Open the database connection
	dbh := GetConnector(dbpath)

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

func WithMutex[T any](fn func(t T) (T, error), arg T) (T, error) {
	mu := sync.Mutex{}
	mu.Lock()
	defer mu.Unlock()

	return fn(arg)
}
