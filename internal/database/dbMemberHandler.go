package database

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"

	"kyri56xcaesar/discord_bots_app/internal/models"

	_ "modernc.org/sqlite"
)

func (dbh *DBHandler) InsertMember(u models.Member) (*models.Member, error) {

	err := u.VerifyMember()
	if err != nil {
		log.Print("Invalid field on Member. ", err.Error())
		return nil, err
	}

	mu := &dbh.mu

	mu.Lock()
	defer mu.Unlock()

	err = dbh.openConnection()
	if err != nil {
		log.Printf("There's been an error getting the DB handler..." + err.Error())
		return nil, err
	}
	defer dbh.DB.Close()

	res, err := dbh.DB.Exec(`INSERT INTO members (guild, username, nickname, avatarurl, displayavatarurl, bannerurl, displaybannerurl, usercolor, joinedat, userstatus, msgcount) 
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		u.Guild, u.Username, u.Nick, u.Avatar, u.DisplayBanner, u.Banner,
		u.DisplayBanner, u.User_color, u.JoinedAt, u.Status, u.MsgCount)
	if err != nil {
		log.Printf("There's been an error inserting the member %v in the DB."+err.Error(), u)
		return nil, err
	}

	lastId, err := res.LastInsertId()
	if err != nil {
		log.Printf("There's been an error retrieving result ID." + err.Error())
		return nil, err
	}

	u.ID = int(lastId)

	// Insert the Roles and Messages now
	var successMCount, successRCount int = 0, 0
	if u.Messages != nil {
		for _, msg := range u.Messages {
			_, err := dbh.DB.Exec(`INSERT INTO messages (userid, content, channel, createdat)
			VALUES (?, ?, ?, ?)`, lastId, msg.Content, msg.Channel, msg.CreatedAt)
			if err != nil {
				log.Printf("Error inserting message %v into the database", msg)
			}
			successMCount++
		}
	}

	if u.Roles != nil {
		for _, role := range u.Roles {
			_, err := dbh.DB.Exec(`INSERT INTO roles (userid, rolename, rolecolor) VALUES (?, ?, ?)`,
				lastId, role.Role_name, role.Role_name)
			if err != nil {
				log.Printf("Error inserting role %v into the database", role)
			}
			successRCount++
		}
	}

	log.Printf("Inserted %d/%d roles, %d/%d messages.", successRCount, len(u.Roles), successMCount, len(u.Messages))

	return &u, err
}

func (dbh *DBHandler) InsertMultipleMembers(members []models.Member) (string, error) {

	mu := &dbh.mu
	mu.Lock()
	defer mu.Unlock()

	err := dbh.openConnection()
	if err != nil {
		log.Printf("There's been an error creating the DB handler: %v", err)
		return "Failed to create DB handler", err
	}
	defer dbh.DB.Close()

	db := dbh.DB

	// Start a transaction
	tx, err := db.Begin()
	if err != nil {
		log.Printf("Failed to begin transaction: %v", err)
		return "Failed to begin transaction", err
	}

	// Prepare the SQL statement for inserting members
	stmt, err := tx.Prepare(`INSERT INTO members (guild, username, nickname, avatarurl, displayavatarurl, bannerurl, displaybannerurl, usercolor, joinedat, userstatus, msgcount) 
        VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`)

	if err != nil {
		log.Printf("Failed to prepare statement: %v", err)
		return "Failed to prepare statement", err
	}
	defer stmt.Close() // Ensure the statement is closed after use

	var successMCount, successRCount, successCount int = 0, 0, 0

	// Loop through each member and try to insert them
	for _, u := range members {
		// Verify the member's data
		err = u.VerifyMember()
		if err != nil {
			log.Printf("Invalid member %+v: %v (Skipping)", u, err.Error())
			continue // Skip faulty member and proceed with the next
		}

		// Execute the insert statement
		res, err := stmt.Exec(u.Guild, u.Username, u.Nick, u.Avatar, u.DisplayAvatar, u.Banner,
			u.DisplayBanner, u.User_color, u.JoinedAt, u.Status, u.MsgCount)
		if err != nil {
			log.Printf("Failed to insert member %v: %v", u, err)
			log.Printf("(Skipping)")
			continue // Skip faulty member and proceed with the next
		}

		lastId, err := res.LastInsertId()
		if err != nil {
			log.Printf("There's been an error retrieving result ID." + err.Error())
			break
		}

		// Increment the success counter if a row was inserted
		successCount++

		if u.Messages != nil {
			for _, msg := range u.Messages {
				_, err := dbh.DB.Exec(`INSERT INTO messages (userid, content, channel, createdat)
			VALUES (?, ?, ?, ?)`, lastId, msg.Content, msg.Channel, msg.CreatedAt)
				if err != nil {
					log.Printf("Error inserting message %v into the database", msg)
				}
				successMCount++
			}
		}

		if u.Roles != nil {
			for _, role := range u.Roles {
				_, err := dbh.DB.Exec(`INSERT INTO roles (userid, rolename, rolecolor) VALUES (?, ?, ?)`,
					lastId, role.Role_name, role.Role_name)
				if err != nil {
					log.Printf("Error inserting role %v into the database", role)
				}
				successRCount++
			}
		}

	}

	// Commit the transaction even if some members were skipped
	if err := tx.Commit(); err != nil {
		log.Printf("Failed to commit transaction: %v", err)
		return "Failed to commit transaction", err
	}

	// Return the number of successful insertions
	return fmt.Sprintf("Successfully inserted %d members", successCount), nil
}

func (dbh *DBHandler) GetAllMembers() ([]*models.Member, error) {

	mu := &dbh.mu

	mu.Lock()
	defer mu.Unlock()

	err := dbh.openConnection()
	if err != nil {
		log.Printf("There's been an error creating the DB handler..." + err.Error())
		return nil, err
	}

	defer dbh.DB.Close()

	rows, err := dbh.DB.Query("SELECT * FROM members")
	if err != nil {
		log.Printf("There is been an error retrieving members from the database." + err.Error())
		return nil, err
	}
	defer rows.Close()

	var members []*models.Member
	memberMap := make(map[int]*models.Member)

	for rows.Next() {

		member := &models.Member{}

		if err := rows.Scan(&member.ID, &member.Guild, &member.Username, &member.Nick,
			&member.Avatar, &member.DisplayAvatar,
			&member.Banner, &member.DisplayBanner,
			&member.User_color, &member.JoinedAt,
			&member.Status, &member.MsgCount); err != nil {
			log.Printf("There's been an error scanning a user from the database." + err.Error())
			return nil, err
		}

		member.Roles = []models.Role{}
		member.Messages = []models.Message{}
		members = append(members, member)
		memberMap[member.ID] = member
	}

	rrows, err := dbh.DB.Query("SELECT * FROM roles")
	if err != nil {
		log.Print("There is been an error retrieving roles with from the database. " + err.Error())
		return nil, err
	}
	defer rrows.Close()

	for rrows.Next() {
		var role models.Role

		if err := rrows.Scan(&role.ID, &role.UID, &role.Role_name, &role.Color); err != nil {
			log.Print("There's been an error scanning user roles from the database " + err.Error())
			return nil, err
		}

		if member, exists := memberMap[role.UID]; exists {
			member.Roles = append(member.Roles, role)
		}

	}

	mrows, err := dbh.DB.Query("SELECT * FROM messages")
	if err != nil {
		log.Print("There is been an error retrieving messages from the database. ", err.Error())
		return nil, err
	}
	defer mrows.Close()

	for mrows.Next() {
		var message models.Message

		if err := mrows.Scan(&message.ID, &message.UID, &message.Content, &message.Channel, &message.CreatedAt); err != nil {
			log.Print("There's been an error scanning user messages from the database " + err.Error())
			return nil, err
		}

		if member, exists := memberMap[message.UID]; exists {
			member.Messages = append(member.Messages, message)
		}

	}

	return members, nil
}

func (dbh *DBHandler) GetMemberByIdentifier(identifier string) (*models.Member, error) {

	mu := &dbh.mu

	mu.Lock()
	defer mu.Unlock()

	err := dbh.openConnection()
	if err != nil {
		log.Printf("There's been an error creating the DB handler..." + err.Error())
		return nil, err
	}

	defer dbh.DB.Close()

	var row *sql.Row

	if isNumeric(identifier) {
		row = dbh.DB.QueryRow("SELECT * FROM members WHERE id = ?", identifier)
	} else {
		row = dbh.DB.QueryRow("SELECT * FROM members WHERE username = ?", identifier)

	}

	member := models.Member{}
	member.Roles = []models.Role{}
	member.Messages = []models.Message{}

	if err := row.Scan(&member.ID, &member.Guild, &member.Username, &member.Nick, &member.Avatar, &member.DisplayAvatar, &member.Banner, &member.DisplayBanner, &member.User_color, &member.JoinedAt, &member.Status, &member.MsgCount); err != nil {
		if err == sql.ErrNoRows {
			log.Printf("Member not found: %v", identifier)
			return nil, nil
		}
		log.Printf("Error scanning member from the row: %v", err)
		return nil, err
	}

	rrows, err := dbh.DB.Query("SELECT * FROM roles WHERE userid = ?", member.ID)
	if err == nil {
		defer rrows.Close()

		for rrows.Next() {
			var role models.Role

			if err := rrows.Scan(&role.ID, &role.UID, &role.Role_name, &role.Color); err != nil {
				log.Printf("There's been an error scanning a role for userid %v %v", member.ID, err.Error())
				break
			}

			member.Roles = append(member.Roles, role)
		}
	}

	mrows, err := dbh.DB.Query("SELECT * FROM messages WHERE userid = ?", member.ID)
	if err == nil {
		defer mrows.Close()

		for mrows.Next() {
			var message models.Message

			if err := mrows.Scan(&message.ID, &message.UID, &message.Content, &message.Channel, &message.CreatedAt); err != nil {
				log.Printf("There's been an error scanning a message for userid %v %v", member.ID, err.Error())
				break
			}

			member.Messages = append(member.Messages, message)
		}
	}

	return &member, nil
}

func (dbh *DBHandler) DeleteMemberByIdentifier(identifier string) (string, error) {

	mu := &dbh.mu

	mu.Lock()
	defer mu.Unlock()

	err := dbh.openConnection()
	if err != nil {
		log.Printf("There's been an error creating the DB handler..." + err.Error())
		return "all not G bro", err
	}

	defer dbh.DB.Close()

	var res sql.Result

	if isNumeric(identifier) {

		res, err = dbh.DB.Exec(`DELETE FROM members WHERE id = ?`, identifier)
	} else {
		res, err = dbh.DB.Exec(`DELETE FROM members WHERE username = ?`, identifier)

	}

	if err != nil {
		log.Printf("There's been an error deleting the member with id %v in the DB."+err.Error(), identifier)
		return "error deleting member", err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		log.Printf("There's been an error retrieving result ID." + err.Error())
		return "error retrieving data", err
	}

	return strconv.FormatInt(rowsAffected, 10), nil
}

func (dbh *DBHandler) UpdateMemberByIdentifier(u models.Member, identifier string) (string, error) {

	mu := &dbh.mu

	mu.Lock()
	defer mu.Unlock()

	err := dbh.openConnection()
	if err != nil {
		log.Printf("There's been an error creating the DB handler..." + err.Error())
		return "all not G bro", err
	}

	defer dbh.DB.Close()

	var res sql.Result

	if isNumeric(identifier) {
		res, err = dbh.DB.Exec(`UPDATE members SET guild = ?, id = ?, username = ?, nickname = ?, avatarurl = ?, 
		displayavatarurl = ?, bannerurl = ?, displaybannerurl = ?, usercolor, 
		joinedat = ?, userstatus = ?, msgcount = ? WHERE id = ?`,
			u.Guild, u.ID, u.Username, u.Nick, u.Avatar, u.DisplayBanner, u.Banner,
			u.DisplayBanner, u.User_color, u.JoinedAt, u.Status, u.MsgCount, identifier)
	} else {
		res, err = dbh.DB.Exec(`UPDATE members SET guild = ?, id = ?, username = ?, nickname = ?, avatarurl = ?, 
		displayavatarurl = ?, bannerurl = ?, displaybannerurl = ?, usercolor, 
		joinedat = ?, userstatus = ?, msgcount = ? WHERE username = ?`,
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
