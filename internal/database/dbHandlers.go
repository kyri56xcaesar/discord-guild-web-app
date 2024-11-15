package database

import (
	"fmt"
	"log"
	"strings"

	"kyri56xcaesar/discord_bots_app/internal/models"
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

func (dbh *DBHandler) Select(tableName string, columns []string, identifiers map[string][]string, limit int, sortField, order string) (interface{}, error) {
	if tableName == "" {
		return nil, fmt.Errorf("no tablename given.")
	}

	var cols string = "*"
	if len(columns) > 0 {
		cols = strings.Join(columns, ", ")
	}

	query := fmt.Sprintf("SELECT %s FROM %s", cols, tableName)

	// Perhaps move to a function
	var whereClauses []string
	var args []interface{}
	i := 1

	if identifiers != nil {
		for key, value := range identifiers {
			whereClauses = append(whereClauses, "(")
			// log.Printf("key: %v, value: %v", key, value)
			for index, param := range value {
				// log.Printf("index: %v, param: %v", index, param)
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

	if len(sortField) > 0 {
		if len(order) <= 0 {
			order = DefaultOrder
		}
		query = fmt.Sprintf("%s ORDER BY %s %s", query, sortField, order)
	}

	if limit > 0 {
		query = fmt.Sprintf("%s LIMIT %d", query, limit)
	}

	// log.Printf("Tablename: %v", tableName)
	// log.Printf("Query ready: %v", query)

	err := dbh.openConnection()
	if err != nil {
		log.Printf("There's been an error creating the DB handler: %v", err)
		return nil, fmt.Errorf("failed to fetch the db handler: %w", err)
	}

	rows, err := dbh.DB.Query(query, args...)
	if err != nil {
		log.Printf("Error executing the Query. :%v", err)
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}

	// Should handle the rows returned.
	// Could be either slice or single of a Member/Bot/Line
	switch tableName {
	case TypeMember:
		var members []models.Member
		for rows.Next() {
			var member models.Member
			err = rows.Scan(member.PtrFieldsSpecific(columns)...)
			if err != nil {
				log.Printf("Error scanning member row: %v", err)
				return nil, fmt.Errorf("failed to scan member row: %w", err)
			}
			members = append(members, member)
		}
		if len(members) == 0 {
			return nil, nil
		}
		return members, nil
	case TypeBot:
		var bots []models.Bot
		for rows.Next() {
			var bot models.Bot
			err = rows.Scan(bot.PtrFieldsSpecific(columns)...)
			if err != nil {
				log.Printf("Error scanning bot row: %v", err)
				return nil, fmt.Errorf("failed to scan bot row: %w", err)
			}
			bots = append(bots, bot)
		}
		if len(bots) == 0 {
			return nil, nil
		}
		return bots, nil
	case TypeLine:
		var lines []models.Line
		for rows.Next() {
			var line models.Line
			err = rows.Scan(line.PtrFieldsSpecific(columns)...)
			if err != nil {
				log.Printf("Error scaning line row: %v", err)
				return nil, fmt.Errorf("failed to scan line row: %w", err)
			}
			lines = append(lines, line)
		}
		if len(lines) == 0 {
			return nil, nil
		}
		return lines, nil
	default:
		// would have exitted earlier if this is the case, during the query forming
		log.Print("I am surprised we reached here.")
		return nil, fmt.Errorf("Invalid table name")
	}
}

func (dbh *DBHandler) Delete(tableName string, identifiers []string) {
}
