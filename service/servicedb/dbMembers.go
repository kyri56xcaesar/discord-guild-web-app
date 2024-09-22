package servicedb

import (
	"log"
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

	db := dbHandler.DB
	defer db.Close()

	res, err := db.Exec(`INSERT INTO members (guild, id, username, nickname, avatarurl, displayavatarurl, bannerurl, displaybannerurl, usercolor, joinedat, userstatus, msgcount) 
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		u.Guild, u.ID, u.User, u.Nick, u.Avatar, u.DisplayBanner, u.Banner,
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

	return "", nil
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

	db := dbHandler.DB
	defer db.Close()

	rows, err := db.Query("SELECT * FROM members")
	if err != nil {
		log.Printf("There is been an error retrieving members from the database." + err.Error())
		return nil, err
	}

	var members []user.User
	for rows.Next() {

		var user user.User
		if err := rows.Scan(&user.ID, &user.Guild, &user.User, &user.Nick, &user.Avatar, &user.DisplayAvatar, &user.Banner, &user.DisplayBanner, &user.User_color, &user.JoinedAt, &user.Status, &user.MsgCount); err != nil {
			log.Printf("There's been an error scanning a user from the database." + err.Error())
			return nil, err
		}

		members = append(members, user)
	}

	return members, nil

}

func GetMemberByID(id int) (*user.User, error) {
	var dBHandler DBHandler

	dbHandler, err := dBHandler.OpenConnection(DBName)
	if err != nil {
		log.Printf("There's been an error creating the DB handler..." + err.Error())
		return nil, err
	}

	mu := &dBHandler.MU

	mu.Lock()
	defer mu.Unlock()

	db := dbHandler.DB
	defer db.Close()

	row := db.QueryRow("SELECT * FROM members WHERE uid = ?", id)
	// if err != nil {
	// 	log.Printf("There is been an error retrieving members from the database." + err.Error())
	// 	return nil, err
	// }

	user := user.User{}

	if err := row.Scan(&user.ID, &user.Guild, &user.User, &user.Nick, &user.Avatar, &user.DisplayAvatar, &user.Banner, &user.DisplayBanner, &user.User_color, &user.JoinedAt, &user.Status, &user.MsgCount); err != nil {
		log.Printf("There's been an error scanning the member from the row." + err.Error())
		return nil, err
	}

	return &user, nil
}

func GetMemberByName(name string) (*user.User, error) {
	var dBHandler DBHandler

	dbHandler, err := dBHandler.OpenConnection(DBName)
	if err != nil {
		log.Printf("There's been an error creating the DB handler..." + err.Error())
		return nil, err
	}

	mu := &dBHandler.MU

	mu.Lock()
	defer mu.Unlock()

	db := dbHandler.DB
	defer db.Close()

	row := db.QueryRow("SELECT * FROM members WHERE username = ?", name)
	// if err != nil {
	// 	log.Printf("There is been an error retrieving members from the database." + err.Error())
	// 	return nil, err
	// }

	user := &user.User{}

	if err := row.Scan(user); err != nil {
		log.Printf("There's been an error scanning the member from the row." + err.Error())
		return nil, err
	}

	return user, nil
}

func DeleteMemberByID(id int) (string, error) {

	var dBHandler DBHandler

	dbHandler, err := dBHandler.OpenConnection(DBName)
	if err != nil {
		log.Printf("There's been an error creating the DB handler..." + err.Error())
		return "all not G bro", err
	}

	mu := &dBHandler.MU

	mu.Lock()
	defer mu.Unlock()

	db := dbHandler.DB
	defer db.Close()

	res, err := db.Exec(`DELETE FROM members WHERE id = ?`, id)

	if err != nil {
		log.Printf("There's been an error deleting the member with id %v in the DB."+err.Error(), id)
		return "error deleting member", err
	}

	lastId, err := res.LastInsertId()
	if err != nil {
		log.Printf("There's been an error retrieving result ID." + err.Error())
		return "error retrieving data", err
	}

	return strconv.FormatInt(lastId, 10), nil

}

func DeleteMemberByName(name string) (string, error) {
	var dBHandler DBHandler

	dbHandler, err := dBHandler.OpenConnection(DBName)
	if err != nil {
		log.Printf("There's been an error creating the DB handler..." + err.Error())
		return "all not G bro", err
	}

	mu := &dBHandler.MU

	mu.Lock()
	defer mu.Unlock()

	db := dbHandler.DB
	defer db.Close()

	res, err := db.Exec(`DELETE FROM members WHERE username = ?`, name)

	if err != nil {
		log.Printf("There's been an error deleting the member with id %v in the DB."+err.Error(), name)
		return "error deleting member", err
	}

	lastId, err := res.LastInsertId()
	if err != nil {
		log.Printf("There's been an error retrieving result ID." + err.Error())
		return "error retrieving data", err
	}

	return strconv.FormatInt(lastId, 10), nil

}

func UpdateMemberByID(u user.User, id int) (string, error) {
	var dBHandler DBHandler

	dbHandler, err := dBHandler.OpenConnection(DBName)
	if err != nil {
		log.Printf("There's been an error creating the DB handler..." + err.Error())
		return "all not G bro", err
	}

	mu := &dBHandler.MU

	mu.Lock()
	defer mu.Unlock()

	db := dbHandler.DB
	defer db.Close()

	res, err := db.Exec(`UPDATE members SET guild = ?, id = ?, username = ?, nickname = ?, avatarurl = ?, 
		displayavatarurl = ?, bannerurl = ?, displaybannerurl = ?, usercolor, 
		joinedat = ?, userstatus = ?, msgcount = ? WHERE id = ?)`,
		u.Guild, u.ID, u.User, u.Nick, u.Avatar, u.DisplayBanner, u.Banner,
		u.DisplayBanner, u.User_color, u.JoinedAt, u.Status, u.MsgCount, id)

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

func UpdateMemberByName(u user.User, name string) (string, error) {
	var dBHandler DBHandler

	dbHandler, err := dBHandler.OpenConnection(DBName)
	if err != nil {
		log.Printf("There's been an error creating the DB handler..." + err.Error())
		return "all not G bro", err
	}

	mu := &dBHandler.MU

	mu.Lock()
	defer mu.Unlock()

	db := dbHandler.DB
	defer db.Close()

	res, err := db.Exec(`UPDATE members SET guild = ?, id = ?, username = ?, nickname = ?, avatarurl = ?, 
		displayavatarurl = ?, bannerurl = ?, displaybannerurl = ?, usercolor, 
		joinedat = ?, userstatus = ?, msgcount = ? WHERE username = ?)`,
		u.Guild, u.ID, u.User, u.Nick, u.Avatar, u.DisplayBanner, u.Banner,
		u.DisplayBanner, u.User_color, u.JoinedAt, u.Status, u.MsgCount, name)

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
