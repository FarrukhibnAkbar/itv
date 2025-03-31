package utils

import (
	"database/sql"
)

// func Hello() {
// 	fmt.Println("Hello from common package")
// }

func NullString(item string) sql.NullString {
	if item != "" {
		return sql.NullString{String: item, Valid: true}
	}
	return sql.NullString{Valid: false}
}

func NullStringToString(ns sql.NullString) string {
    if ns.Valid {
        return ns.String
    }
    return "" 
}