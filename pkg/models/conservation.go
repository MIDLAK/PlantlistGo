package models

/* Мера сохранения */
type Conservation struct {
	ID			int
	Name		string
	Description	string
	Start		string	//начало принятия мер
	End			string	//конец принятия мер
}