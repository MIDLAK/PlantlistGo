package models

import (
	"database/sql"
)

/* Публикация */
type Publication struct {
	ID			int
	Link		sql.NullString
	Biblio		string			//библиографическое описание
	Description	sql.NullString
}