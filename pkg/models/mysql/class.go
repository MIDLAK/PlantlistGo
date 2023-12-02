package mysql

import (
	"database/sql"
	"GoTest/pkg/models"
	"errors"
)

//добавление нового класса в базу данных
func (plm *PlantlistModel) InsertClass(class string) (int64, error) {
	result, err := plm.DB.Exec(`INSERT IGNORE INTO class (id, name) VALUES (NULL, ?);`, class)
	if err != nil {
		return 0, err
	}

	//id добавленного класса
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}

//получение данных класса по его названию
func (plm *PlantlistModel) GetClass(className string) (*models.Class, error) {
	selectClass := `SELECT id, name FROM class WHERE name = ?;`
	row := plm.DB.QueryRow(selectClass, className)

	class := &models.Class{}
	err := row.Scan(&class.ID, &class.Name)
	if err != nil {
		//если модели в программе и БД не совпадают
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}
	
	return class, nil
}

//обновление класса для растения
func (plm *PlantlistModel) UpdateClass(className string) error {
	_, err := plm.DB.Exec(`UPDATE class SET name=? WHERE name=?;`, className)
	if err != nil {
		return err
	}

	return nil
}