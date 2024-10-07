package database

import (
	"fmt"
	"log"
	"strconv"

	"kyri56xcaesar/discord_bots_app/internal/models"
)

// Bots
func (dbh *DBHandler) GetAllBots() ([]models.Bot, error) {

	mu := &dbh.MU

	mu.Lock()
	defer mu.Unlock()

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

	var bots []models.Bot
	for rows.Next() {

		var bot models.Bot

		if err := rows.Scan(&bot.ID, &bot.Guild, &bot.Name, &bot.Avatar,
			&bot.Banner, &bot.CreatedAt, &bot.Author, &bot.Status, &bot.IsSinger,
		); err != nil {
			log.Printf("There's been an error scanning a user from the database." + err.Error())
			return nil, err
		}

		bots = append(bots, bot)
	}

	return bots, nil
}

func (dbh *DBHandler) InsertMutipleBots(bots []models.Bot) (string, error) {

	mu := &dbh.MU
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
	stmt, err := tx.Prepare(`INSERT INTO bots (botguild, botname, avatarurl, bannerurl, createdat, author, botstatus, isSinger) 
        VALUES (?, ?, ?, ?, ?, ?, ?, ?)`)

	if err != nil {
		log.Printf("Failed to prepare statement: %v", err)
		return "Failed to prepare statement", err
	}
	defer stmt.Close() // Ensure the statement is closed after use

	// Initialize counter for successful insertions
	successCount := 0

	// Loop through each member and try to insert them
	for _, b := range bots {
		// Verify the member's data
		err = b.VerifyBot()
		if err != nil {
			log.Printf("Invalid bot %+v: %v (Skipping)", b, err.Error())
			continue // Skip faulty member and proceed with the next
		}

		// Execute the insert statement
		_, err := stmt.Exec(b.Guild, b.Name, b.Avatar, b.Banner, b.CreatedAt, b.Author,
			b.Status, b.IsSinger)
		if err != nil {
			log.Printf("Failed to insert bot %v: %v", b, err)
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
	return fmt.Sprintf("Successfully inserted %d bots", successCount), nil
}

func (dbh *DBHandler) GetBotByIdentifier(identifier string) (*models.Bot, error) {
	mu := &dbh.MU
	mu.Lock()
	defer mu.Unlock()

	err := dbh.openConnection()
	if err != nil {
		log.Print("There's been an error getting the DB handler! ", err.Error())
		return nil, err
	}
	defer dbh.DB.Close()

	row := dbh.DB.QueryRow("SELECT * FROM bots WHERE botid = ?", identifier)

	bot := models.Bot{}

	if err := row.Scan(&bot.ID, &bot.Guild, &bot.Name, &bot.Avatar,
		&bot.Banner, &bot.CreatedAt, &bot.Author, &bot.Status, &bot.IsSinger,
	); err != nil {
		log.Printf("There's been an error scanning the Line from rows. %v", err.Error())
		return nil, err
	}

	return &bot, nil
}

func (dbh *DBHandler) InsertBot(b *models.Bot) (string, error) {

	err := b.VerifyBot()
	if err != nil {
		log.Print("Invalid field on Bot. ", err.Error())
		return "0", err
	}

	mu := &dbh.MU

	mu.Lock()
	defer mu.Unlock()

	err = dbh.openConnection()
	if err != nil {
		log.Printf("There's been an error getting the DB handler..." + err.Error())
		return "bad db", err
	}
	defer dbh.DB.Close()

	res, err := dbh.DB.Exec(`INSERT INTO bots (botguild, botname, avatarurl, bannerurl, createdat, author, botstatus, isSinger) 
        VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		b.Guild, b.Name, b.Avatar, b.Banner, b.CreatedAt, b.Author, b.Status, b.IsSinger)

	if err != nil {
		log.Printf("There's been an error inserting the bot %v in the DB."+err.Error(), b)
		return "error inserting bot", err
	}

	lastId, err := res.LastInsertId()
	if err != nil {
		log.Printf("There's been an error retrieving result ID." + err.Error())
		return "error retrieving data", err
	}

	return fmt.Sprintf("{'status':%v}", strconv.FormatInt(lastId, 10)), err
}

func (dbh *DBHandler) UpdateBotByIdentifier(b *models.Bot, identifier string) (string, error) {
	mu := &dbh.MU

	mu.Lock()
	defer mu.Unlock()

	err := dbh.openConnection()
	if err != nil {
		log.Printf("There's been an error getting the DB handler..." + err.Error())
		return "all not G bro", err
	}

	defer dbh.DB.Close()

	res, err := dbh.DB.Exec(`UPDATE bots SET botguild = ?, botname = ?, avatarurl = ?, bannerurl = ?, createdat = ?, author = ?, botstatus = ?, isSinger = ? WHERE botid = ?`,
		b.Guild, b.Name, b.Avatar, b.Banner, b.CreatedAt, b.Author, b.Status, b.IsSinger, identifier)

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
	mu := &dbh.MU

	mu.Lock()
	defer mu.Unlock()

	err := dbh.openConnection()
	if err != nil {
		log.Printf("There's been an error getting the DB handler..." + err.Error())
		return "Nothing deleted", err
	}

	defer dbh.DB.Close()

	res, err := dbh.DB.Exec(`DELETE FROM bots WHERE botid = ?`, identifier)
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

// Lines
func (dbh *DBHandler) GetBotLines() ([]models.Line, error) {

	mu := &dbh.MU
	mu.Lock()
	defer mu.Unlock()

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

		if err := rows.Scan(&line.ID, &line.BID, &line.Phrase, &line.Author, &line.To, &line.LineType, &line.CreatedAt); err != nil {
			log.Printf("There's been an error scanning the Line from rows. %v", err.Error())
			return nil, err
		}

		lines = append(lines, line)
	}

	return lines, nil
}

func (dbh *DBHandler) InsertMultipleLines(lines []models.Line) (string, error) {

	mu := &dbh.MU
	mu.Lock()
	defer mu.Unlock()

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
		_, err := stmt.Exec(l.BID, l.Phrase, l.Author, l.To, l.LineType)
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

	mu := &dbh.MU
	mu.Lock()
	defer mu.Unlock()

	err := dbh.openConnection()
	if err != nil {
		log.Print("There's been an error getting the DB handler! ", err.Error())
		return nil, err
	}
	defer dbh.DB.Close()

	row := dbh.DB.QueryRow("SELECT * FROM lines WHERE id = ?", identifier)

	line := models.Line{}

	if err := row.Scan(&line.ID, &line.BID, &line.Phrase, &line.Author, &line.To, &line.LineType, &line.CreatedAt); err != nil {
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

	mu := &dbh.MU

	mu.Lock()
	defer mu.Unlock()

	err = dbh.openConnection()
	if err != nil {
		log.Printf("There's been an error getting the DB handler..." + err.Error())
		return "bad db", err
	}
	defer dbh.DB.Close()

	res, err := dbh.DB.Exec(`INSERT INTO lines (bid, phrase, author, toid, ltype)
		VALUES (?, ?, ?, ?, ?)`,
		l.BID, l.Phrase, l.Author, l.To, l.LineType)

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

func (dbh *DBHandler) UpdateLineByIndentifier(l models.Line, identifier string) (string, error) {
	mu := &dbh.MU

	mu.Lock()
	defer mu.Unlock()

	err := dbh.openConnection()
	if err != nil {
		log.Printf("There's been an error getting the DB handler..." + err.Error())
		return "all not G bro", err
	}

	defer dbh.DB.Close()

	res, err := dbh.DB.Exec(`UPDATE lines SET lineid = ?, bid = ?, phrase = ?, author = ?, toid = ?, ltype = ?, createdat = ? WHERE id = ?`,
		l.ID, l.BID, l.Phrase, l.Author, l.To, l.LineType, l.CreatedAt, identifier)

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

func (dbh *DBHandler) DeleteLineByIndentifier(identifier string) (string, error) {

	mu := &dbh.MU

	mu.Lock()
	defer mu.Unlock()

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
