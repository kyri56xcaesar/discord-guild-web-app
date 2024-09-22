package servicedb

import (
	"database/sql"
	"log"
	"regexp"
	"strconv"

	"kyri56xcaesar/discord_bots_app/guild/user"
)

const DBName string = "servicedb/myapp.db"

func InsertMember(u user.User) (string, error) {

	var dBHandler DBHandler

	dbHandler, err := dBHandler.OpenConnection(DBName)
	if err != nil {
		log.Printf("There's been an error creating the DB handler..." + err.Error())
		return "all not G bro", err
	}

	mu := &dBHandler.MU

	mu.Lock()
	defer mu.Unlock()

	defer dbHandler.DB.Close()

	res, err := dbHandler.DB.Exec(`INSERT INTO members (guild, id, username, nickname, avatarurl, displayavatarurl, bannerurl, displaybannerurl, usercolor, joinedat, userstatus, msgcount) 
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		u.Guild, u.ID, u.Username, u.Nick, u.Avatar, u.DisplayBanner, u.Banner,
		u.DisplayBanner, u.User_color, u.JoinedAt, u.Status, u.MsgCount)

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

func InsertMultipleMembers(members []user.User) (string, error) {
	var dBHandler DBHandler

	dbHandler, err := dBHandler.OpenConnection(DBName)
	if err != nil {
		log.Printf("There's been an error creating the DB handler: %v", err)
		return "Failed to create DB handler", err
	}
	defer dbHandler.DB.Close()

	mu := &dbHandler.MU
	mu.Lock()
	defer mu.Unlock()

	db := dbHandler.DB

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
		_, err := stmt.Exec(u.Guild, u.ID, u.Username, u.Nick, u.Avatar, u.DisplayAvatar, u.Banner,
			u.DisplayBanner, u.User_color, u.JoinedAt, u.Status, u.MsgCount)
		if err != nil {
			log.Printf("Failed to insert member %v: %v", u, err)
			// Rollback the transaction if there's an error
			tx.Rollback()
			return "Error inserting member", err
		}
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		log.Printf("Failed to commit transaction: %v", err)
		return "Failed to commit transaction", err
	}

	return "Members inserted successfully", nil
}

func GetAllMembers() ([]user.User, error) {

	var dbHandler DBHandler

	_, err := dbHandler.OpenConnection(DBName)
	if err != nil {
		log.Printf("There's been an error creating the DB handler..." + err.Error())
		return nil, err
	}

	mu := &dbHandler.MU

	mu.Lock()
	defer mu.Unlock()

	defer dbHandler.DB.Close()

	rows, err := dbHandler.DB.Query("SELECT * FROM members")
	if err != nil {
		log.Printf("There is been an error retrieving members from the database." + err.Error())
		return nil, err
	}

	var members []user.User
	for rows.Next() {

		var user user.User
		if err := rows.Scan(&user.ID, &user.Guild, &user.Username, &user.Nick,
			&user.Avatar, &user.DisplayAvatar,
			&user.Banner, &user.DisplayBanner,
			&user.User_color, &user.JoinedAt,
			&user.Status, &user.MsgCount); err != nil {
			log.Printf("There's been an error scanning a user from the database." + err.Error())
			return nil, err
		}

		members = append(members, user)
	}

	return members, nil

}

func isNumeric(s string) bool {
	re := regexp.MustCompile(`^[0-9]+$`)
	return re.MatchString(s)
}

func GetMemberByIdentifier(identifier string) (*user.User, error) {
	var dBHandler DBHandler

	dbHandler, err := dBHandler.OpenConnection(DBName)
	if err != nil {
		log.Printf("There's been an error creating the DB handler..." + err.Error())
		return nil, err
	}

	mu := &dBHandler.MU

	mu.Lock()
	defer mu.Unlock()

	defer dbHandler.DB.Close()

	var row *sql.Row

	if isNumeric(identifier) {
		row = dbHandler.DB.QueryRow("SELECT * FROM members WHERE id = ?", identifier)
	} else {
		row = dbHandler.DB.QueryRow("SELECT * FROM members WHERE username = ?", identifier)

	}
	// if err != nil {
	// 	log.Printf("There is been an error retrieving members from the database." + err.Error())
	// 	return nil, err
	// }

	user := user.User{}

	if err := row.Scan(&user.ID, &user.Guild, &user.Username, &user.Nick, &user.Avatar, &user.DisplayAvatar, &user.Banner, &user.DisplayBanner, &user.User_color, &user.JoinedAt, &user.Status, &user.MsgCount); err != nil {
		log.Printf("There's been an error scanning the member from the row." + err.Error())
		return nil, err
	}

	return &user, nil
}

func DeleteMemberByIdentifier(identifier string) (string, error) {

	var dBHandler DBHandler

	dbHandler, err := dBHandler.OpenConnection(DBName)
	if err != nil {
		log.Printf("There's been an error creating the DB handler..." + err.Error())
		return "all not G bro", err
	}

	mu := &dBHandler.MU

	mu.Lock()
	defer mu.Unlock()

	defer dbHandler.DB.Close()

	var res sql.Result

	if isNumeric(identifier) {

		res, err = dbHandler.DB.Exec(`DELETE FROM members WHERE id = ?`, identifier)
	} else {
		res, err = dbHandler.DB.Exec(`DELETE FROM members WHERE username = ?`, identifier)

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

func UpdateMemberByIdentifier(u user.User, identifier string) (string, error) {
	var dBHandler DBHandler

	dbHandler, err := dBHandler.OpenConnection(DBName)
	if err != nil {
		log.Printf("There's been an error creating the DB handler..." + err.Error())
		return "all not G bro", err
	}

	mu := &dBHandler.MU

	mu.Lock()
	defer mu.Unlock()

	defer dbHandler.DB.Close()

	var res sql.Result

	if isNumeric(identifier) {
		res, err = dbHandler.DB.Exec(`UPDATE members SET guild = ?, id = ?, username = ?, nickname = ?, avatarurl = ?, 
		displayavatarurl = ?, bannerurl = ?, displaybannerurl = ?, usercolor, 
		joinedat = ?, userstatus = ?, msgcount = ? WHERE id = ?)`,
			u.Guild, u.ID, u.Username, u.Nick, u.Avatar, u.DisplayBanner, u.Banner,
			u.DisplayBanner, u.User_color, u.JoinedAt, u.Status, u.MsgCount, identifier)
	} else {
		res, err = dbHandler.DB.Exec(`UPDATE members SET guild = ?, id = ?, username = ?, nickname = ?, avatarurl = ?, 
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
