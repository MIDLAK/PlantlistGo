package models

/* Растение */
type Plant struct {
	ID				int
	RusName			string
	LatinName		string
	Description		string
	Image			string	//путь к файлу
}