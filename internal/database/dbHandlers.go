package database

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	"kyri56xcaesar/discord_bots_app/internal/models"
)

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

func (dbh *DBHandler) Search(tableName string) {
}

func (dbh *DBHandler) Select(columns, tableName string, identifiers map[string][]string, relational bool) (*sql.Rows, error) {
	if tableName == "" {
		return nil, fmt.Errorf("No tablename given.")
	}

	cols := "*"
	if len(columns) > 0 {
		cols = columns
	}

	query := fmt.Sprintf("SELECT %s FROM %s", cols, tableName)

	var whereClauses []string
	var args []interface{}
	i := 1
	for key, value := range identifiers {
		whereClauses = append(whereClauses, "(")
		for index, param := range value {
			if index < len(value)-1 {
				whereClauses = append(whereClauses, fmt.Sprintf("%s = ? OR ", key))
			} else {
				whereClauses = append(whereClauses, fmt.Sprintf("%s = ?", key))
			}
			args = append(args, param)
		}
		if len(identifiers) > i {
			whereClauses = append(whereClauses, ") AND ")
		} else {
			whereClauses = append(whereClauses, ")")
		}
		i++
	}

	query = fmt.Sprintf("%s WHERE %s", query, strings.Join(whereClauses, ""))

	log.Print(query)

	rows, err := dbh.DB.Query(query, args...)
	if err != nil {
		log.Printf("Failed to execute query")
		return nil, err
	}

	for rows.Next() {
		member := models.Member{}

		err := rows.Scan(&member)
		if err != nil {
			log.Printf("error scanning result")
			return nil, err
		} else {
			log.Print(member)
		}
	}

	// Members, Bots or Lines

	return nil, nil
}

func (dbh *DBHandler) Delete(tableName string, identifiers []string) {
}
