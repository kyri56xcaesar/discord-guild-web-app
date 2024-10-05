package database

import (
	"fmt"
	"kyri56xcaesar/discord_bots_app/internal/models"
	"log"
)

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
			log.Printf("Invalid member %+v: %v (Skipping)", b, err.Error())
			continue // Skip faulty member and proceed with the next
		}

		// Execute the insert statement
		_, err := stmt.Exec(b.Guild, b.Name, b.Avatar, b.Banner, b.CreatedAt, b.Author,
			b.Status, b.IsSinger)
		if err != nil {
			log.Printf("Failed to insert member %v: %v", b, err)
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
	return fmt.Sprintf("Successfully inserted %d members", successCount), nil
}
