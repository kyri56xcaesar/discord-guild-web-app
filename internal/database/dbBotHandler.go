package database

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"strings"

	"kyri56xcaesar/discord-guild-web-app/internal/models"
	"kyri56xcaesar/discord-guild-web-app/internal/utils"
)

// Bots
func (dbh *DBHandler) GetAllBots() ([]*models.Bot, error) {
	err := dbh.openConnection()
	if err != nil {
		log.Printf("There's been an error creating the DB handler..." + err.Error())
		return nil, err
	}

	defer dbh.DB.Close()

	rows, err := dbh.DB.Query("SELECT * FROM bots")
	if err != nil {
		log.Printf("There is been an error retrieving members from the database." + err.Error())
		return nil, err
	}

	defer rows.Close()

	var bots []*models.Bot
	botMap := make(map[int]*models.Bot)

	for rows.Next() {

		bot := &models.Bot{}

		if err := rows.Scan(bot.PtrFieldsDB()...); err != nil {
			log.Printf("There's been an error scanning a user from the database." + err.Error())
			return nil, err
		}

		bot.Lines = []models.Line{}
		bots = append(bots, bot)
		botMap[bot.Id] = bot
	}

	lrows, err := dbh.DB.Query("SELECT * FROM lines")
	if err != nil {
		log.Printf("There's been an error retrieving lines from database, %v", err.Error())
		return bots, err
	}
	defer lrows.Close()

	for lrows.Next() {
		var line models.Line

		if err := lrows.Scan(line.PtrFieldsDB()...); err != nil {
			log.Print("There's been an error scanning line. " + err.Error())
			return bots, err
		}

		if bot, exists := botMap[line.Bid]; exists {
			bot.Lines = append(bot.Lines, line)
		}
	}

	return bots, nil
}

func (dbh *DBHandler) GetMultipleBotsByIdentifiers(identifiers []string) ([]*models.Bot, error) {
	err := dbh.openConnection()
	if err != nil {
		log.Printf("There's been an error creating the DB handler..." + err.Error())
		return nil, err
	}
	defer dbh.DB.Close()

	query := "SELECT * FROM bots WHERE id IN (?" + strings.Repeat(",?", len(identifiers)-1) + ")"

	// Execute the query with the provided identifiers
	rows, err := dbh.DB.Query(query, utils.InterfaceSlice(identifiers))
	if err != nil {
		log.Printf("Error retrieving bots from the database: %v", err)
		return nil, err
	}
	defer rows.Close()

	var bots []*models.Bot
	botMap := make(map[int]*models.Bot)

	for rows.Next() {

		bot := &models.Bot{}

		if err := rows.Scan(bot.PtrFieldsDB()...); err != nil {
			log.Printf("There's been an error scanning a member from the database." + err.Error())
			return nil, err
		}

		bot.Lines = []models.Line{}
		bots = append(bots, bot)
		botMap[bot.Id] = bot
	}

	lrows, err := dbh.DB.Query("SELECT * FROM lines")
	if err != nil {
		log.Printf("There's been an error retrieving lines from database, %v", err.Error())
		return bots, err
	}
	defer lrows.Close()

	for lrows.Next() {
		var line models.Line

		if err := lrows.Scan(line.PtrFieldsDB()...); err != nil {
			log.Print("There's been an error scanning line. " + err.Error())
			return bots, err
		}

		if bot, exists := botMap[line.Bid]; exists {
			bot.Lines = append(bot.Lines, line)
		}
	}

	return bots, nil
}

func (dbh *DBHandler) InsertMultipleBots(bots []models.Bot) (string, error) {
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
	stmt, err := tx.Prepare(`INSERT INTO bots (guild, username, avatarurl, bannerurl, createdat, author, status, issinger) 
        VALUES (?, ?, ?, ?, ?, ?, ?, ?)`)
	if err != nil {
		log.Printf("Failed to prepare statement: %v", err)
		return "Failed to prepare statement", err
	}
	defer stmt.Close() // Ensure the statement is closed after use

	// Initialize counter for successful insertions
	var successCount, successLCount int

	for _, b := range bots {
		// Verify the member's data
		err = b.VerifyBot()
		if err != nil {
			log.Printf("Invalid bot %+v: %v (Skipping)", b, err.Error())
			continue // Skip faulty member and proceed with the next
		}

		// Execute the insert statement
		res, err := stmt.Exec(b.Guild, b.Username, b.Avatarurl, b.Bannerurl, b.Createdat, b.Author, b.Status, b.Issinger)
		if err != nil {
			log.Printf("Failed to insert bot %v: %v", b, err)
			log.Printf("(Skipping)")
			continue // Skip faulty member and proceed with the next
		}

		lastId, err := res.LastInsertId()
		if err != nil {
			log.Printf("There's been an error retrieving result ID. " + err.Error())
			break
		}
		// Increment the success counter if a row was inserted
		successCount++

		if b.Lines != nil {
			for _, line := range b.Lines {
				_, err := dbh.DB.Exec(`
				INSERT INTO lines 
					(bid, phrase, author, toid, ltype, createdat)
				VALUES 
					(?, ?, ?, ?, ?, ?)`,
					lastId, line.Phrase, line.Author, line.Toid, line.Ltype, line.Createdat)
				if err != nil {
					log.Printf("Error inserting line %v into the database, %v", line, err.Error())
				}
				successLCount++
			}
		}

	}

	// Commit the transaction even if some members were skipped
	if err := tx.Commit(); err != nil {
		log.Printf("Failed to commit transaction: %v", err)
		return "Failed to commit transaction", err
	}

	// Return the number of successful insertions
	return fmt.Sprintf("Successfully inserted %d bots", successCount), nil
}

func (dbh *DBHandler) GetBotIdentifiers(identifier string) ([]string, error) {
	err := dbh.openConnection()
	if err != nil {
		log.Print("There's been an error getting the DB handler! ", err.Error())
		return nil, err
	}
	defer dbh.DB.Close()

	if !AllowedBotCols[identifier] {
		return nil, fmt.Errorf("invalid identifier name: %s", identifier)
	}

	ident := identifier[:len(identifier)-1]

	query := fmt.Sprintf("SELECT %s FROM bots", ident)

	rows, err := dbh.DB.Query(query)
	if err != nil {
		log.Print("There's been an error retrieving bot data. " + err.Error())
		return nil, err
	}

	var results []string

	for rows.Next() {
		var content string

		if err := rows.Scan(&content); err != nil {
			log.Printf("There's been an error scanning the bot data. %v", err)
			return nil, err
		}

		results = append(results, content)
	}
	return results, nil
}

func (dbh *DBHandler) GetBotByIdentifier(identifier string) (*models.Bot, error) {
	err := dbh.openConnection()
	if err != nil {
		log.Print("There's been an error getting the DB handler! ", err.Error())
		return nil, err
	}
	defer dbh.DB.Close()

	row := dbh.DB.QueryRow("SELECT * FROM bots WHERE id = ?", identifier)

	bot := models.Bot{}
	bot.Lines = []models.Line{}

	if err := row.Scan(bot.PtrFieldsDB()...); err != nil {
		log.Printf("There's been an error scanning the Line from rows. %v", err.Error())
		return nil, err
	}

	lrows, err := dbh.DB.Query("SELECT * FROM lines WHERE bid = ?", bot.Id)
	if err == nil {
		defer lrows.Close()

		for lrows.Next() {
			var line models.Line

			if err := lrows.Scan(line.PtrFieldsDB()...); err != nil {
				log.Printf("There's been an error scanning the line for botid %v %v", bot.Id, err.Error())
				break
			}

			bot.Lines = append(bot.Lines, line)
		}
	}

	return &bot, nil
}

func (dbh *DBHandler) InsertBot(b *models.Bot) (string, error) {
	err := b.VerifyBot()
	if err != nil {
		log.Print("Invalid field on Bot. ", err.Error())
		return "0", err
	}

	err = dbh.openConnection()
	if err != nil {
		log.Printf("There's been an error getting the DB handler..." + err.Error())
		return "bad db", err
	}
	defer dbh.DB.Close()

	res, err := dbh.DB.Exec(`INSERT INTO bots (guild, username, avatarurl, bannerurl, createdat, author, status, issinger) 
        VALUES (?, ?, ?, ?, ?, ?, ?, ?)`, b.Guild, b.Username, b.Avatarurl, b.Bannerurl, b.Createdat, b.Author, b.Status, b.Issinger)
	if err != nil {
		log.Printf("There's been an error inserting the bot %v in the DB."+err.Error(), b)
		return "error inserting bot", err
	}

	lastId, err := res.LastInsertId()
	if err != nil {
		log.Printf("There's been an error retrieving result ID." + err.Error())
		return "error retrieving data", err
	}

	var successLCount int
	if b.Lines != nil {
		for _, line := range b.Lines {
			_, err := dbh.DB.Exec(`
				INSERT INTO lines 
					(bid, phrase, author, toid, ltype, createdat)
				VALUES 
					(?, ?, ?, ?, ?, ?)`, line.Bid, line.Phrase, line.Author, line.Toid, line.Ltype, line.Createdat)
			if err != nil {
				log.Printf("Error inserting line %v into the database", line)
			}
			successLCount++

		}
	}

	log.Printf("Inserted: %d/%d lines", successLCount, len(b.Lines))

	return fmt.Sprintf("{'status':%v}", strconv.FormatInt(lastId, 10)), err
}

func (dbh *DBHandler) UpdateBotByIdentifier(b *models.Bot, identifier string) (string, error) {
	err := dbh.openConnection()
	if err != nil {
		log.Printf("There's been an error getting the DB handler..." + err.Error())
		return "all not G bro", err
	}

	defer dbh.DB.Close()

	res, err := dbh.DB.Exec(`UPDATE bots SET guild = ?, username = ?, avatarurl = ?, bannerurl = ?, createdat = ?, author = ?, status = ?, issinger = ? WHERE id = ?`, b.Guild,
		b.Username,
		b.Avatarurl,
		b.Bannerurl,
		b.Createdat,
		b.Author,
		b.Status,
		b.Issinger,
		identifier)
	if err != nil {
		log.Printf("There's been an error updating the bot %v in the DB."+err.Error(), b)
		return "error inserting bot", err
	}

	lastId, err := res.LastInsertId()
	if err != nil {
		log.Printf("There's been an error retrieving result ID." + err.Error())
		return "error retrieving data", err
	}

	return strconv.FormatInt(lastId, 10), nil
}

func (dbh *DBHandler) DeleteBotByIdentifier(identifier string) (string, error) {
	err := dbh.openConnection()
	if err != nil {
		log.Printf("There's been an error getting the DB handler..." + err.Error())
		return "Nothing deleted", err
	}

	defer dbh.DB.Close()

	res, err := dbh.DB.Exec(`DELETE FROM bots WHERE id = ?`, identifier)
	if err != nil {
		log.Printf("There's been an error deleting the bot with id %v in the DB."+err.Error(), identifier)
		return "error deleting bot", err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		log.Printf("There's been an error retrieving result ID." + err.Error())
		return "error retrieving data", err
	}

	return strconv.FormatInt(rowsAffected, 10), nil
}

func (dbh *DBHandler) DeleteMultipleBotsByIdentifiers(identifiers []string) (string, error) {
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

	queryID := `DELETE FROM bots WHERE id = ?`
	queryUsername := `DELETE FROM bots WHERE username = ?`

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
			log.Printf("Error deleting bot with identifier %v: %v", identifier, err)
			return "Error deleting one or more bots", err
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

	return fmt.Sprintf("Successfully deleted %d bots", totalDeleted), nil
}

func (dbh *DBHandler) GetLineIdentifiers(identifier string) ([]string, error) {
	err := dbh.openConnection()
	if err != nil {
		log.Print("There's been an error getting the DB handler! ", err.Error())
		return nil, err
	}
	defer dbh.DB.Close()

	if !AllowedLineCols[identifier] {
		return nil, fmt.Errorf("invalid identifier name: %s", identifier)
	}

	ident := identifier[:len(identifier)-1]
	query := fmt.Sprintf("SELECT %s FROM lines", ident)

	rows, err := dbh.DB.Query(query)
	if err != nil {
		log.Print("There's been an error retrieving line data. " + err.Error())
		return nil, err
	}

	var results []string

	for rows.Next() {
		var content string

		if err := rows.Scan(&content); err != nil {
			log.Printf("There's been an error scanning the line data. %v", err)
			return nil, err
		}

		results = append(results, content)
	}
	return results, nil
}

// Lines
func (dbh *DBHandler) GetBotLines() ([]models.Line, error) {
	err := dbh.openConnection()
	if err != nil {
		log.Print("There's been an error getting the DB handler! ", err.Error())
		return nil, err
	}
	defer dbh.DB.Close()

	rows, err := dbh.DB.Query("SELECT * FROM lines")
	if err != nil {
		log.Printf("There is been an error retrieving members from the database." + err.Error())
		return nil, err
	}

	var lines []models.Line

	for rows.Next() {
		var line models.Line

		if err := rows.Scan(line.PtrFieldsDB()...); err != nil {
			log.Printf("There's been an error scanning the Line from rows. %v", err.Error())
			return nil, err
		}

		lines = append(lines, line)
	}

	return lines, nil
}

func (dbh *DBHandler) GetMultipleLinesByIdentifiers(identifiers []string) ([]models.Line, error) {
	err := dbh.openConnection()
	if err != nil {
		log.Print("There's been an error getting the DB handler! ", err.Error())
		return nil, err
	}
	defer dbh.DB.Close()

	query := "SELECT * FROM lines WHERE id IN (?" + strings.Repeat(",?", len(identifiers)-1) + ")"

	// Execute the query with the provided identifiers
	rows, err := dbh.DB.Query(query, utils.InterfaceSlice(identifiers))
	if err != nil {
		log.Printf("Error retrieving lines from the database: %v", err)
		return nil, err
	}
	defer rows.Close()

	var lines []models.Line

	for rows.Next() {
		var line models.Line

		if err := rows.Scan(line.PtrFieldsDB()...); err != nil {
			log.Printf("There's been an error scanning the Line from rows. %v", err.Error())
			return nil, err
		}

		lines = append(lines, line)
	}

	return lines, nil
}

func (dbh *DBHandler) InsertMultipleLines(lines []models.Line) (string, error) {
	err := dbh.openConnection()
	if err != nil {
		log.Printf("There's been an error getting the DB handler: %v", err)
		return "Failed to get the DB handler", err
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
	stmt, err := tx.Prepare(`INSERT INTO lines (bid, phrase, author, toid, ltype)
		VALUES (?, ?, ?, ?, ?)`)
	if err != nil {
		log.Printf("Failed to prepare statement: %v", err)
		return "Failed to prepare statement", err
	}
	defer stmt.Close() // Ensure the statement is closed after use

	// Initialize counter for successful insertions
	successCount := 0

	// Loop through each member and try to insert them
	for _, l := range lines {
		// Verify the member's data
		err = l.VerifyLine()
		if err != nil {
			log.Printf("Invalid line %+v: %v (Skipping)", l, err.Error())
			continue // Skip faulty member and proceed with the next
		}

		// Execute the insert statement
		_, err := stmt.Exec(l.Bid, l.Phrase, l.Author, l.Toid, l.Ltype)
		if err != nil {
			log.Printf("Failed to insert line %v: %v", l, err)
			log.Printf("(Skipping)")
			continue // Skip faulty member and proceed with the next
		}

		// Increment the success counter if a row was inserted
		successCount++

	}

	// Commit the transaction even if some members were skipped
	if err := tx.Commit(); err != nil {
		log.Printf("Failed to commit transaction: %v", err)
		return "Failed to commit transaction", err
	}

	// Return the number of successful insertions
	return fmt.Sprintf("Successfully inserted %d lines", successCount), nil
}

func (dbh *DBHandler) GetBotLineByIdentifier(identifier string) (*models.Line, error) {
	err := dbh.openConnection()
	if err != nil {
		log.Print("There's been an error getting the DB handler! ", err.Error())
		return nil, err
	}
	defer dbh.DB.Close()

	row := dbh.DB.QueryRow("SELECT * FROM lines WHERE id = ?", identifier)

	line := models.Line{}

	if err := row.Scan(line.PtrFieldsDB()...); err != nil {
		log.Printf("There's been an error scanning the Line from rows. %v", err.Error())
		return nil, err
	}

	return &line, nil
}

func (dbh *DBHandler) InsertLine(l *models.Line) (string, error) {
	err := l.VerifyLine()
	if err != nil {
		log.Print("Invalid field on Line. ", err.Error())
		return "0", err
	}

	err = dbh.openConnection()
	if err != nil {
		log.Printf("There's been an error getting the DB handler..." + err.Error())
		return "bad db", err
	}
	defer dbh.DB.Close()

	res, err := dbh.DB.Exec(`INSERT INTO lines (bid, phrase, author, toid, ltype)
		VALUES (?, ?, ?, ?, ?)`,
		l.Bid, l.Phrase, l.Author, l.Toid, l.Ltype)
	if err != nil {
		log.Printf("There's been an error inserting the line %v in the DB."+err.Error(), l)
		return "error inserting line", err
	}

	lastId, err := res.LastInsertId()
	if err != nil {
		log.Printf("There's been an error retrieving result ID." + err.Error())
		return "error retrieving data", err
	}

	return fmt.Sprintf("{'status':%v}", strconv.FormatInt(lastId, 10)), err
}

func (dbh *DBHandler) UpdateLineByIdentifier(l models.Line, identifier string) (string, error) {
	err := dbh.openConnection()
	if err != nil {
		log.Printf("There's been an error getting the DB handler..." + err.Error())
		return "all not G bro", err
	}

	defer dbh.DB.Close()

	res, err := dbh.DB.Exec(`UPDATE lines SET id = ?, bid = ?, phrase = ?, author = ?, toid = ?, ltype = ?, createdat = ? WHERE id = ?`,
		l.Id, l.Bid, l.Phrase, l.Author, l.Toid, l.Ltype, l.Createdat, identifier)
	if err != nil {
		log.Printf("There's been an error updating the line %v in the DB."+err.Error(), l)
		return "error inserting line", err
	}

	lastId, err := res.LastInsertId()
	if err != nil {
		log.Printf("There's been an error retrieving result ID." + err.Error())
		return "error retrieving data", err
	}

	return strconv.FormatInt(lastId, 10), nil
}

func (dbh *DBHandler) DeleteLineByIdentifier(identifier string) (string, error) {
	err := dbh.openConnection()
	if err != nil {
		log.Printf("There's been an error getting the DB handler..." + err.Error())
		return "Nothing deleted", err
	}

	defer dbh.DB.Close()

	res, err := dbh.DB.Exec(`DELETE FROM lines WHERE id = ?`, identifier)
	if err != nil {
		log.Printf("There's been an error deleting the line with id %v in the DB."+err.Error(), identifier)
		return "error deleting line", err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		log.Printf("There's been an error retrieving result ID." + err.Error())
		return "error retrieving data", err
	}

	return strconv.FormatInt(rowsAffected, 10), nil
}

func (dbh *DBHandler) DeleteMultipleLinesByIdentifiers(identifiers []string) (string, error) {
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

	queryID := `DELETE FROM lines WHERE id = ?`

	totalDeleted := 0
	for _, identifier := range identifiers {
		var res sql.Result
		if utils.IsNumeric(identifier) {
			res, err = tx.Exec(queryID, identifier)
		} else {
			continue
		}

		if err != nil {
			tx.Rollback()
			log.Printf("Error deleting line with identifier %v: %v", identifier, err)
			return "Error deleting one or more lines", err
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

	return fmt.Sprintf("Successfully deleted %d lines", totalDeleted), nil
}
