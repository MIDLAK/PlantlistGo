package mysql

import (
	"database/sql"
	"GoTest/pkg/models"
	"errors"
)

//добавление нового рода в базу данных
func (plm *PlantlistModel) InsertGenus(genus string) (int64, error) {
	result, err := plm.DB.Exec(`INSERT IGNORE INTO genus (id, name) VALUES (NULL, ?);`, genus)
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

//получение данных рода по его названию
func (plm *PlantlistModel) GetGenus(genName string) (*models.Genus, error) {
	selectGenus := `SELECT id, name FROM genus WHERE name = ?;`
	row := plm.DB.QueryRow(selectGenus, genName)

	gen := &models.Genus{}
	err := row.Scan(&gen.ID, &gen.Name)
	if err != nil {
		//если модели в программе и БД не совпадают
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}
	
	return gen, nil
}

//обновление данных по семейству таксона
func (plm *PlantlistModel) UpdateGenus(genName string) error {
	_, err := plm.DB.Exec(`UPDATE genus SET name=? WHERE name=?;`, genName, genName)
	if err != nil {
		return err
	}

	return nil
}