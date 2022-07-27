package repository

import "database/sql"

func NewNullString(strPtr *string) sql.NullString {
	if strPtr != nil {
		return sql.NullString{String: *strPtr, Valid: true}
	}
	return sql.NullString{}
}
