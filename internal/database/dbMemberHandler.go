package database

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"strings"

	"kyri56xcaesar/discord-guild-web-app/internal/models"
	"kyri56xcaesar/discord-guild-web-app/internal/utils"

	_ "modernc.org/sqlite"
)

// Members
func (dbh *DBHandler) InsertMember(u models.Member) (*models.Member, error) {
	err := u.VerifyMember()
	if err != nil {
		log.Print("Invalid field on Member. ", err.Error())
		return nil, err
	}

	err = dbh.openConnection()
	if err != nil {
		log.Printf("There's been an error getting the DB handler..." + err.Error())
		return nil, err
	}
	defer dbh.DB.Close()

	res, err := dbh.DB.Exec(`INSERT INTO members (guild, username, nickname, leaguename, avatarurl, displayavatarurl, bannerurl, displaybannerurl, usercolor, joinedat, status, msgcount) 
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		u.Guild, u.Username, u.Nickname, u.Avatarurl, u.Displaybannerurl, u.Bannerurl,
		u.Displaybannerurl, u.Usercolor, u.Joinedat, u.Status, u.Msgcount)
	if err != nil {
		log.Printf("There's been an error inserting the member %v in the DB."+err.Error(), u)
		return nil, err
	}

	lastId, err := res.LastInsertId()
	if err != nil {
		log.Printf("There's been an error retrieving result ID." + err.Error())
		return nil, err
	}

	u.Id = int(lastId)

	// Insert the Roles and Messages now
	var successMCount, successRCount int
	if u.Usermessages != nil {
		for _, msg := range u.Usermessages {
			_, err := dbh.DB.Exec(`INSERT INTO messages (userid, content, channel, createdat)
			VALUES (?, ?, ?, ?)`, lastId, msg.Content, msg.Channel, msg.Createdat)
			if err != nil {
				log.Printf("Error inserting message %v into the database", msg)
			}
			successMCount++
		}
	}

	if u.Userroles != nil {
		for _, role := range u.Userroles {
			_, err := dbh.DB.Exec(`INSERT INTO roles (userid, rolename, rolecolor) VALUES (?, ?, ?)`,
				lastId, role.Rolename, role.Rolename)
			if err != nil {
				log.Printf("Error inserting role %v into the database", role)
			}
			successRCount++
		}
	}

	log.Printf("Inserted %d/%d roles, %d/%d messages.", successRCount, len(u.Userroles), successMCount, len(u.Usermessages))

	return &u, err
}

func (dbh *DBHandler) InsertMultipleMembers(members []models.Member) ([]models.Member, error) {
	err := dbh.openConnection()
	if err != nil {
		log.Printf("There's been an error creating the DB handler: %v", err)
		return nil, err
	}
	defer dbh.DB.Close()
	db := dbh.DB

	// Start a transaction
	tx, err := db.Begin()
	if err != nil {
		log.Printf("Failed to begin transaction: %v", err)
		return nil, err
	}

	// Prepare the SQL statement for inserting members
	stmt, err := tx.Prepare(`INSERT INTO members (guild, username, nickname, leaguename, avatarurl, displayavatarurl, bannerurl, displaybannerurl, usercolor, joinedat, status, msgcount) 
        VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`)
	if err != nil {
		log.Printf("Failed to prepare statement: %v", err)
		return nil, err
	}
	defer stmt.Close() // Ensure the statement is closed after use

	var successMCount, successRCount, successCount int

	// Loop through each member and try to insert them
	for _, u := range members {
		// Verify the member's data
		err = u.VerifyMember()
		if err != nil {
			log.Printf("Invalid member %+v: %v (Skipping)", u, err.Error())
			continue // Skip faulty member and proceed with the next
		}

		// Execute the insert statement
		res, err := stmt.Exec(u.Guild, u.Username, u.Nickname, u.Avatarurl, u.Displayavatarurl, u.Bannerurl,
			u.Displaybannerurl, u.Usercolor, u.Joinedat, u.Status, u.Msgcount)
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

		if u.Usermessages != nil {
			for _, msg := range u.Usermessages {
				_, err := dbh.DB.Exec(`INSERT INTO messages (userid, content, channel, createdat)
			VALUES (?, ?, ?, ?)`, lastId, msg.Content, msg.Channel, msg.Createdat)
				if err != nil {
					log.Printf("Error inserting message %v into the database", msg)
				}
				successMCount++
			}
		}

		if u.Userroles != nil {
			for _, role := range u.Userroles {
				_, err := dbh.DB.Exec(`INSERT INTO roles (userid, rolename, rolecolor) VALUES (?, ?, ?)`,
					lastId, role.Rolename, role.Rolename)
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
		return nil, err
	}

	if successCount == 0 {
		return nil, fmt.Errorf("Failed to insert members %v", members)
	}

	// Return the number of successful insertions
	return members, nil
}

func (dbh *DBHandler) GetAllMembers() ([]*models.Member, error) {
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

		if err := rows.Scan(member.PtrFieldsDB()...); err != nil {
			log.Printf("There's been an error scanning a user from the database." + err.Error())
			return nil, err
		}

		member.Userroles = []models.Role{}
		member.Usermessages = []models.Message{}
		members = append(members, member)
		memberMap[member.Id] = member
	}

	rrows, err := dbh.DB.Query("SELECT * FROM roles")
	if err != nil {
		log.Print("There is been an error retrieving roles with from the database. " + err.Error())
		return nil, err
	}
	defer rrows.Close()

	for rrows.Next() {
		var role models.Role

		if err := rrows.Scan(&role.Id, &role.Userid, &role.Rolename, &role.Rolecolor); err != nil {
			log.Print("There's been an error scanning user roles from the database " + err.Error())
			return nil, err
		}

		if member, exists := memberMap[role.Userid]; exists {
			member.Userroles = append(member.Userroles, role)
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

		if err := mrows.Scan(message.PtrFieldsDB()...); err != nil {
			log.Print("There's been an error scanning user messages from the database " + err.Error())
			return nil, err
		}

		if member, exists := memberMap[message.Userid]; exists {
			member.Usermessages = append(member.Usermessages, message)
		}

	}

	return members, nil
}

func (dbh *DBHandler) GetMultipleMembersByIdentifiers(identifiers []string) ([]*models.Member, error) {
	err := dbh.openConnection()
	if err != nil {
		log.Printf("There's been an error creating the DB handler..." + err.Error())
		return nil, err
	}
	defer dbh.DB.Close()

	// queryUsername := "SELECT * FROM members WHERE username IN "
	query := "SELECT * FROM members WHERE id IN (?" + strings.Repeat(",?", len(identifiers)-1) + ")"

	// Execute the query with the provided identifiers
	rows, err := dbh.DB.Query(query, utils.InterfaceSlice(identifiers))
	if err != nil {
		log.Printf("Error retrieving members from the database: %v", err)
		return nil, err
	}
	defer rows.Close()

	var members []*models.Member
	memberMap := make(map[int]*models.Member)

	for rows.Next() {

		member := &models.Member{}

		if err := rows.Scan(member.PtrFieldsDB()...); err != nil {
			log.Printf("There's been an error scanning a user from the database." + err.Error())
			return nil, err
		}

		member.Userroles = []models.Role{}
		member.Usermessages = []models.Message{}
		members = append(members, member)
		memberMap[member.Id] = member
	}

	rrows, err := dbh.DB.Query("SELECT * FROM roles")
	if err != nil {
		log.Print("There is been an error retrieving roles with from the database. " + err.Error())
		return nil, err
	}
	defer rrows.Close()

	for rrows.Next() {
		var role models.Role

		if err := rrows.Scan(&role.Id, &role.Userid, &role.Rolename, &role.Rolecolor); err != nil {
			log.Print("There's been an error scanning user roles from the database " + err.Error())
			return nil, err
		}

		if member, exists := memberMap[role.Userid]; exists {
			member.Userroles = append(member.Userroles, role)
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

		if err := mrows.Scan(message.PtrFieldsDB()...); err != nil {
			log.Print("There's been an error scanning user messages from the database " + err.Error())
			return nil, err
		}

		if member, exists := memberMap[message.Userid]; exists {
			member.Usermessages = append(member.Usermessages, message)
		}

	}

	return members, nil
}

func (dbh *DBHandler) GetMemberByIdentifier(identifier string) (*models.Member, error) {
	err := dbh.openConnection()
	if err != nil {
		log.Printf("There's been an error creating the DB handler..." + err.Error())
		return nil, err
	}
	defer dbh.DB.Close()

	var row *sql.Row

	if utils.IsNumeric(identifier) {
		row = dbh.DB.QueryRow("SELECT * FROM members WHERE id = ?", identifier)
	} else {
		row = dbh.DB.QueryRow("SELECT * FROM members WHERE username = ?", identifier)
	}

	member := models.Member{}
	member.Userroles = []models.Role{}
	member.Usermessages = []models.Message{}

	if err := row.Scan(member.PtrFieldsDB()...); err != nil {
		if err == sql.ErrNoRows {
			log.Printf("Member not found: %v", identifier)
			return nil, nil
		}
		log.Printf("Error scanning member from the row: %v", err)
		return nil, err
	}

	rrows, err := dbh.DB.Query("SELECT * FROM roles WHERE userid = ?", member.Id)
	if err == nil {
		defer rrows.Close()

		for rrows.Next() {
			var role models.Role

			if err := rrows.Scan(&role.Id, &role.Userid, &role.Rolename, &role.Rolecolor); err != nil {
				log.Printf("There's been an error scanning a role for userid %v %v", member.Id, err.Error())
				break
			}

			member.Userroles = append(member.Userroles, role)
		}
	}

	mrows, err := dbh.DB.Query("SELECT * FROM messages WHERE userid = ?", member.Id)
	if err == nil {
		defer mrows.Close()

		for mrows.Next() {
			var message models.Message

			if err := mrows.Scan(message.PtrFieldsDB()...); err != nil {
				log.Printf("There's been an error scanning a message for userid %v %v", member.Id, err.Error())
				break
			}

			member.Usermessages = append(member.Usermessages, message)
		}
	}

	return &member, nil
}

func (dbh *DBHandler) DeleteMemberByIdentifier(identifier string) (string, error) {
	err := dbh.openConnection()
	if err != nil {
		log.Printf("There's been an error creating the DB handler..." + err.Error())
		return "all not G bro", err
	}

	defer dbh.DB.Close()

	var res sql.Result

	if utils.IsNumeric(identifier) {
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

func (dbh *DBHandler) DeleteMultipleMembersByIdentifiers(identifiers []string) (string, error) {
	err := dbh.openConnection()
	if err != nil {
		log.Printf("There's been an error creating the DB handler... %v", err)
		return "Error opening DB connection", err
	}
	defer dbh.DB.Close()

	tx, err := dbh.DB.Begin()
	if err != nil {
		log.Printf("Error starting transaction: %v", err)
		return "Error starting transaction", err
	}

	queryID := `DELETE FROM members WHERE id = ?`
	queryUsername := `DELETE FROM members WHERE username = ?`

	totalDeleted := 0
	for _, identifier := range identifiers {
		var res sql.Result
		if utils.IsNumeric(identifier) {
			res, err = tx.Exec(queryID, identifier)
		} else {
			res, err = tx.Exec(queryUsername, identifier)
		}

		if err != nil {
			tx.Rollback()
			log.Printf("Error deleting member with identifier %v: %v", identifier, err)
			return "Error deleting one or more members", err
		}

		rowsAffected, err := res.RowsAffected()
		if err != nil {
			tx.Rollback()
			log.Printf("Error retrieving affected rows for identifier %v: %v", identifier, err)
			return "Error retrieving affected rows", err
		}

		totalDeleted += int(rowsAffected)
	}

	if err := tx.Commit(); err != nil {
		log.Printf("Error committing transaction: %v", err)
		return "Error committing transaction", err
	}

	return fmt.Sprintf("Successfully deleted %d members", totalDeleted), nil
}

func (dbh *DBHandler) UpdateMemberByIdentifier(u models.Member, identifier string) (string, error) {
	err := dbh.openConnection()
	if err != nil {
		log.Printf("There's been an error creating the DB handler..." + err.Error())
		return "all not G bro", err
	}

	defer dbh.DB.Close()

	var res sql.Result

	if utils.IsNumeric(identifier) {
		res, err = dbh.DB.Exec(`UPDATE members SET guild = ?, id = ?, username = ?, nickname = ?, leaguename = ?, avatarurl = ?, 
		displayavatarurl = ?, bannerurl = ?, displaybannerurl = ?, usercolor = ?, 
		joinedat = ?, status = ?, msgcount = ? WHERE id = ?`,
			u.Guild, u.Id, u.Username, u.Nickname, u.Leaguename, u.Avatarurl, u.Displaybannerurl, u.Bannerurl,
			u.Displaybannerurl, u.Usercolor, u.Joinedat, u.Status, u.Msgcount, identifier)
	} else {
		res, err = dbh.DB.Exec(`UPDATE members SET guild = ?, id = ?, username = ?, nickname = ?, leaguename = ?, avatarurl = ?, 
		displayavatarurl = ?, bannerurl = ?, displaybannerurl = ?, usercolor, 
		joinedat = ?, status = ?, msgcount = ? WHERE username = ?`,
			u.Guild, u.Id, u.Username, u.Nickname, u.Leaguename, u.Avatarurl, u.Displaybannerurl, u.Bannerurl,
			u.Displaybannerurl, u.Usercolor, u.Joinedat, u.Status, u.Msgcount, identifier)
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

func (dbh *DBHandler) InsertRole(r models.Role) (*models.Role, error) {
	err := r.VerifyRole()
	if err != nil {
		log.Print("Invalid field on Role. ", err.Error())
		return nil, err
	}
	err = dbh.openConnection()
	if err != nil {
		log.Printf("There's been an error getting the DB handler..." + err.Error())
		return nil, err
	}
	defer dbh.DB.Close()

	res, err := dbh.DB.Exec(`INSERT INTO roles (userid, rolename, rolecolor) 
		VALUES (?, ?, ?)`,
		r.Userid, r.Rolename, r.Rolecolor)
	if err != nil {
		log.Printf("There's been an error inserting the role %v in the DB."+err.Error(), r)
		return nil, err
	}

	lastId, err := res.LastInsertId()
	if err != nil {
		log.Printf("There's been an error retrieving result ID." + err.Error())
		return nil, err
	}

	r.Id = int(lastId)

	return &r, err
}

func (dbh *DBHandler) InsertMultipleRoles(roles []models.Role) (string, error) {
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
	stmt, err := tx.Prepare(`INSERT INTO roles (userid, rolename, rolecolor) 
		VALUES (?, ?, ?)`)
	if err != nil {
		log.Printf("Failed to prepare statement: %v", err)
		return "Failed to prepare statement", err
	}
	defer stmt.Close() // Ensure the statement is closed after use

	successCount := 0

	// Loop through each member and try to insert them
	for _, r := range roles {
		// Verify the member's data
		err = r.VerifyRole()
		if err != nil {
			log.Printf("Invalid role %+v: %v (Skipping)", r, err.Error())
			continue // Skip faulty member and proceed with the next
		}

		// Execute the insert statement
		res, err := stmt.Exec(r.Userid, r.Rolename, r.Rolecolor)
		if err != nil {
			log.Printf("Failed to insert role %v: %v", r, err)
			log.Printf("(Skipping)")
			continue // Skip faulty member and proceed with the next
		}

		lastId, err := res.LastInsertId()
		if err != nil {
			log.Printf("There's been an error retrieving result ID." + err.Error())
			break
		}

		r.Id = int(lastId)

		// Increment the success counter if a row was inserted
		successCount++

	}

	// Commit the transaction even if some members were skipped
	if err := tx.Commit(); err != nil {
		log.Printf("Failed to commit transaction: %v", err)
		return "Failed to commit transaction", err
	}

	// Return the number of successful insertions
	return fmt.Sprintf("Successfully inserted %d roles", successCount), nil
}

func (dbh *DBHandler) GetAllRoles() ([]*models.Role, error) {
	err := dbh.openConnection()
	if err != nil {
		log.Printf("There's been an error creating the DB handler..." + err.Error())
		return nil, err
	}

	defer dbh.DB.Close()

	rows, err := dbh.DB.Query("SELECT * FROM roles")
	if err != nil {
		log.Printf("There is been an error retrieving roles from the database." + err.Error())
		return nil, err
	}
	defer rows.Close()

	var roles []*models.Role

	for rows.Next() {

		role := &models.Role{}

		if err := rows.Scan(&role.Id, &role.Userid, &role.Rolename, &role.Rolecolor); err != nil {
			log.Printf("There's been an error scanning a role from the database." + err.Error())
			return nil, err
		}

		roles = append(roles, role)
	}

	return roles, nil
}

func (dbh *DBHandler) GetRoleByIdentifier(identifier string) (*models.Role, error) {
	err := dbh.openConnection()
	if err != nil {
		log.Printf("There's been an error creating the DB handler..." + err.Error())
		return nil, err
	}

	defer dbh.DB.Close()

	var row *sql.Row

	if utils.IsNumeric(identifier) {
		row = dbh.DB.QueryRow("SELECT * FROM roles WHERE id = ?", identifier)
	} else {
		row = dbh.DB.QueryRow("SELECT * FROM roles WHERE rolename = ?", identifier)
	}

	role := models.Role{}

	if err := row.Scan(&role.Id, &role.Userid, &role.Rolename, &role.Rolecolor); err != nil {
		if err == sql.ErrNoRows {
			log.Printf("Role not found: %v", identifier)
			return nil, nil
		}
		log.Printf("Error scanning role from the row: %v", err)
		return nil, err
	}

	return &role, nil
}

func (dbh *DBHandler) DeleteRoleByIdentifier(identifier string) (string, error) {
	err := dbh.openConnection()
	if err != nil {
		log.Printf("There's been an error creating the DB handler..." + err.Error())
		return "all not G bro", err
	}

	defer dbh.DB.Close()

	var res sql.Result

	if utils.IsNumeric(identifier) {
		res, err = dbh.DB.Exec(`DELETE FROM roles WHERE id = ?`, identifier)
	} else {
		res, err = dbh.DB.Exec(`DELETE FROM roles WHERE rolename = ?`, identifier)
	}

	if err != nil {
		log.Printf("There's been an error deleting the role with id %v in the DB. "+err.Error(), identifier)
		return "error deleting role", err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		log.Printf("There's been an error retrieving result ID." + err.Error())
		return "error retrieving data", err
	}

	return strconv.FormatInt(rowsAffected, 10), nil
}

func (dbh *DBHandler) UpdateRoleByIdentifier(r models.Role, identifier string) (string, error) {
	err := dbh.openConnection()
	if err != nil {
		log.Printf("There's been an error creating the DB handler..." + err.Error())
		return "all not G bro", err
	}

	defer dbh.DB.Close()

	var res sql.Result

	if utils.IsNumeric(identifier) {
		res, err = dbh.DB.Exec(`UPDATE role SET userid = ?, rolename = ?, rolecolor = ? WHERE id = ?`,
			r.Userid, r.Rolename, r.Rolecolor, identifier)
	} else {
		res, err = dbh.DB.Exec(`UPDATE role SET userid = ?, rolename = ?, rolecolor = ? WHERE rolename = ?`,
			r.Userid, r.Rolename, r.Rolecolor, identifier)
	}

	if err != nil {
		log.Printf("There's been an error updating the role %v in the DB."+err.Error(), r)
		return "error updating role", err
	}

	lastId, err := res.LastInsertId()
	if err != nil {
		log.Printf("There's been an error retrieving result ID." + err.Error())
		return "error retrieving data", err
	}

	return strconv.FormatInt(lastId, 10), nil
}
