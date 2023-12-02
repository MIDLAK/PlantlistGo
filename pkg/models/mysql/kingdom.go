package mysql

import (
	"database/sql"
	"GoTest/pkg/models"
	"errors"
)

//добавление нового царства в базу данных
func (plm *PlantlistModel) InsertKingdom(kingdom string) (int64, error) {
	result, err := plm.DB.Exec(`INSERT IGNORE INTO kingdom (id, name) VALUES (NULL, ?);`, kingdom)
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

//получение данных царства по его названию
func (plm *PlantlistModel) GetKingdom(kingdomName string) (*models.Kingdom, error) {
	selectKingdom := `SELECT id, name FROM kingdom WHERE name = ?;`
	row := plm.DB.QueryRow(selectKingdom, kingdomName)

	kingdom := &models.Kingdom{}
	err := row.Scan(&kingdom.ID, &kingdom.Name)
	if err != nil {
		//если модели в программе и БД не совпадают
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}
	
	return kingdom, nil
}

//обновление данных по семейству таксона
func (plm *PlantlistModel) UpdateKingdom(kinName string) error {
	_, err := plm.DB.Exec(`UPDATE kingdom SET name=? WHERE name=?;`, kinName, kinName)
	if err != nil {
		return err
	}

	return nil
}