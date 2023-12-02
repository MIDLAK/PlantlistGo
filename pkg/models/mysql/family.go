package mysql

import (
	"database/sql"
	"GoTest/pkg/models"
	"errors"
)

//добавление нового семейства в базу данных
func (plm *PlantlistModel) InsertFamily(family string) (int64, error) {
	result, err := plm.DB.Exec(`INSERT IGNORE INTO family (id, name) VALUES (NULL, ?);`, family)
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

//получение данных семейства по его названию
func (plm *PlantlistModel) GetFamily(fmlName string) (*models.Family, error) {
	selectFamily := `SELECT id, name FROM family WHERE name = ?;`
	row := plm.DB.QueryRow(selectFamily, fmlName)

	family := &models.Family{}
	err := row.Scan(&family.ID, &family.Name)
	if err != nil {
		//если модели в программе и БД не совпадают
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}
	
	return family, nil
}

//обновление данных по семейству таксона
func (plm *PlantlistModel) UpdateFamily(famName string) error {
	_, err := plm.DB.Exec(`UPDATE family SET name=? WHERE name=?;`, famName, famName)
	if err != nil {
		return err
	}

	return nil
}