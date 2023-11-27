package models

/* Публикация */
type Publication struct {
	ID			int
	Link		string
	Biblio		string	//библиографическое описание
	Description	string
}