package models

import (
	"time"
)

/* Мера сохранения */
type Conservation struct {
	ID			int
	Name		string
	Description	string
	Start		time.Time	//начало принятия мер
	End			time.Time	//конец принятия мер
}