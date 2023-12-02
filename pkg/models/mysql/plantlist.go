package mysql

import (
	"database/sql"
	"errors"
)

//обёртка пула подключений
type PlantlistModel struct {
	DB *sql.DB
}

var ErrNoRecord = errors.New("Подходящей записи не найдено")