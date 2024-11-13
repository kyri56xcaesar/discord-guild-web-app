package database

import (
	"database/sql"
	"fmt"
	"log"
	"strings"
)

func (dbh *DBHandler) Metrics(mtype string) (string, error) {
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

func (dbh *DBHandler) Select(tableName, columns string, identifiers map[string][]string, limit int, sortField, order string) (*sql.Rows, error) {
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

	if identifiers != nil {
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

	}
	log.Print(query)

	// Members, Bots or Lines

	return nil, nil
}

func (dbh *DBHandler) Delete(tableName string, identifiers []string) {
}
