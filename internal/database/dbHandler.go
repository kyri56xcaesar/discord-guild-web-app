package database

import (
	"database/sql"
	"log"
	"os"
	"regexp"
	"strconv"
	"sync"

	"kyri56xcaesar/discord_bots_app/internal/models"

	_ "modernc.org/sqlite"
)

const (
	DBName            string = "myapp.db"
	initSQLScriptPath string = "C:\\Users\\kyria\\Documents\\Coding\\Discord Bots\\internal\\database\\db_init.sql"
)

var (
	dBHandler DBHandler
	ID        int = 0
)

type DBHandler struct {
	DB *sql.DB
	MU sync.Mutex

	dbfile string
}

// Opens a connection to the database and holds reference to the Struct Handler
// Should Close The Connection!
func (dbH *DBHandler) openConnection() error {
	db, err := sql.Open("sqlite", dbH.dbfile)
	if err != nil {
		return err
	}

	dbH.DB = db

	return nil
}

// sqlite path db file
func InitDB(dbpath string) error {
	// Create and init the database!
	dBHandler.dbfile = dbpath
	err := dBHandler.openConnection()
	if err != nil {
		log.Print("Error initializing database connection..., will continue in mem: " + err.Error())
	}

	defer dBHandler.DB.Close()

	fileContent, err := os.ReadFile(initSQLScriptPath)
	if err != nil {
		log.Printf("There was an error reading the sql file, will use a default instead...: " + err.Error())
		dBHandler.RunSQLscript(INITsql)

	} else {
		dBHandler.RunSQLscript(string(fileContent))
	}

	return err

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
		);`

func InsertMember(u models.Member) (string, error) {

	mu := &dBHandler.MU

	mu.Lock()
	defer mu.Unlock()

	err := dBHandler.openConnection()
	if err != nil {
		log.Printf("There's been an error creating the DB handler..." + err.Error())
		return "all not G bro", err
	}

	defer dBHandler.DB.Close()

	res, err := dBHandler.DB.Exec(`INSERT INTO members (guild, id, username, nickname, avatarurl, displayavatarurl, bannerurl, displaybannerurl, usercolor, joinedat, userstatus, msgcount) 
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		u.Guild, ID, u.Username, u.Nick, u.Avatar, u.DisplayBanner, u.Banner,
		u.DisplayBanner, u.User_color, u.JoinedAt, u.Status, u.MsgCount)
	ID++

	if err != nil {
		log.Printf("There's been an error inserting the member %v in the DB."+err.Error(), u)
		return "error inserting member", err
	}

	lastId, err := res.LastInsertId()
	if err != nil {
		log.Printf("There's been an error retrieving result ID." + err.Error())
		return "error retrieving data", err
	}

	return strconv.FormatInt(lastId, 10), nil
}

func InsertMultipleMembers(members []models.Member) (string, error) {

	mu := &dBHandler.MU
	mu.Lock()
	defer mu.Unlock()

	err := dBHandler.openConnection()
	if err != nil {
		log.Printf("There's been an error creating the DB handler: %v", err)
		return "Failed to create DB handler", err
	}
	defer dBHandler.DB.Close()

	db := dBHandler.DB

	// Start a transaction
	tx, err := db.Begin()
	if err != nil {
		log.Printf("Failed to begin transaction: %v", err)
		return "Failed to begin transaction", err
	}

	// Prepare the SQL statement for inserting members
	stmt, err := tx.Prepare(`INSERT INTO members (guild, id, username, nickname, avatarurl, displayavatarurl, bannerurl, displaybannerurl, usercolor, joinedat, userstatus, msgcount) 
        VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`)
	if err != nil {
		log.Printf("Failed to prepare statement: %v", err)
		return "Failed to prepare statement", err
	}
	defer stmt.Close() // Ensure the statement is closed after use

	for _, u := range members {
		_, err := stmt.Exec(u.Guild, ID, u.Username, u.Nick, u.Avatar, u.DisplayAvatar, u.Banner,
			u.DisplayBanner, u.User_color, u.JoinedAt, u.Status, u.MsgCount)
		if err != nil {
			log.Printf("Failed to insert member %v: %v", u, err)
			// Rollback the transaction if there's an error
			tx.Rollback()
			return "Error inserting member", err
		}
		ID++
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		log.Printf("Failed to commit transaction: %v", err)
		return "Failed to commit transaction", err
	}

	return "Members inserted successfully", nil
}

func GetAllMembers() ([]models.Member, error) {

	mu := &dBHandler.MU

	mu.Lock()
	defer mu.Unlock()

	err := dBHandler.openConnection()
	if err != nil {
		log.Printf("There's been an error creating the DB handler..." + err.Error())
		return nil, err
	}

	defer dBHandler.DB.Close()

	rows, err := dBHandler.DB.Query("SELECT * FROM members")
	if err != nil {
		log.Printf("There is been an error retrieving members from the database." + err.Error())
		return nil, err
	}

	var members []models.Member
	for rows.Next() {

		var member models.Member

		if err := rows.Scan(&member.ID, &member.Guild, &member.Username, &member.Nick,
			&member.Avatar, &member.DisplayAvatar,
			&member.Banner, &member.DisplayBanner,
			&member.User_color, &member.JoinedAt,
			&member.Status, &member.MsgCount); err != nil {
			log.Printf("There's been an error scanning a user from the database." + err.Error())
			return nil, err
		}

		members = append(members, member)
	}

	return members, nil

}

func isNumeric(s string) bool {
	re := regexp.MustCompile(`^[0-9]+$`)
	return re.MatchString(s)
}

func GetMemberByIdentifier(identifier string) (*models.Member, error) {

	mu := &dBHandler.MU

	mu.Lock()
	defer mu.Unlock()

	err := dBHandler.openConnection()
	if err != nil {
		log.Printf("There's been an error creating the DB handler..." + err.Error())
		return nil, err
	}

	defer dBHandler.DB.Close()

	var row *sql.Row

	if isNumeric(identifier) {
		row = dBHandler.DB.QueryRow("SELECT * FROM members WHERE id = ?", identifier)
	} else {
		row = dBHandler.DB.QueryRow("SELECT * FROM members WHERE username = ?", identifier)

	}
	// if err != nil {
	// 	log.Printf("There is been an error retrieving members from the database." + err.Error())
	// 	return nil, err
	// }

	member := models.Member{}

	if err := row.Scan(&member.ID, &member.Guild, &member.Username, &member.Nick, &member.Avatar, &member.DisplayAvatar, &member.Banner, &member.DisplayBanner, &member.User_color, &member.JoinedAt, &member.Status, &member.MsgCount); err != nil {
		log.Printf("There's been an error scanning the member from the row." + err.Error())
		return nil, err
	}

	return &member, nil
}

func DeleteMemberByIdentifier(identifier string) (string, error) {

	mu := &dBHandler.MU

	mu.Lock()
	defer mu.Unlock()

	err := dBHandler.openConnection()
	if err != nil {
		log.Printf("There's been an error creating the DB handler..." + err.Error())
		return "all not G bro", err
	}

	defer dBHandler.DB.Close()

	var res sql.Result

	if isNumeric(identifier) {

		res, err = dBHandler.DB.Exec(`DELETE FROM members WHERE id = ?`, identifier)
	} else {
		res, err = dBHandler.DB.Exec(`DELETE FROM members WHERE username = ?`, identifier)

	}

	if err != nil {
		log.Printf("There's been an error deleting the member with id %v in the DB."+err.Error(), identifier)
		return "error deleting member", err
	}

	lastId, err := res.LastInsertId()
	if err != nil {
		log.Printf("There's been an error retrieving result ID." + err.Error())
		return "error retrieving data", err
	}

	return strconv.FormatInt(lastId, 10), nil

}

func UpdateMemberByIdentifier(u models.Member, identifier string) (string, error) {

	mu := &dBHandler.MU

	mu.Lock()
	defer mu.Unlock()

	err := dBHandler.openConnection()
	if err != nil {
		log.Printf("There's been an error creating the DB handler..." + err.Error())
		return "all not G bro", err
	}

	defer dBHandler.DB.Close()

	var res sql.Result

	if isNumeric(identifier) {
		res, err = dBHandler.DB.Exec(`UPDATE members SET guild = ?, id = ?, username = ?, nickname = ?, avatarurl = ?, 
		displayavatarurl = ?, bannerurl = ?, displaybannerurl = ?, usercolor, 
		joinedat = ?, userstatus = ?, msgcount = ? WHERE id = ?)`,
			u.Guild, u.ID, u.Username, u.Nick, u.Avatar, u.DisplayBanner, u.Banner,
			u.DisplayBanner, u.User_color, u.JoinedAt, u.Status, u.MsgCount, identifier)
	} else {
		res, err = dBHandler.DB.Exec(`UPDATE members SET guild = ?, id = ?, username = ?, nickname = ?, avatarurl = ?, 
		displayavatarurl = ?, bannerurl = ?, displaybannerurl = ?, usercolor, 
		joinedat = ?, userstatus = ?, msgcount = ? WHERE username = ?)`,
			u.Guild, u.ID, u.Username, u.Nick, u.Avatar, u.DisplayBanner, u.Banner,
			u.DisplayBanner, u.User_color, u.JoinedAt, u.Status, u.MsgCount, identifier)

	}

	if err != nil {
		log.Printf("There's been an error updating the member %v in the DB."+err.Error(), u)
		return "error inserting member", err
	}

	lastId, err := res.LastInsertId()
	if err != nil {
		log.Printf("There's been an error retrieving result ID." + err.Error())
		return "error retrieving data", err
	}

	return strconv.FormatInt(lastId, 10), nil
}
