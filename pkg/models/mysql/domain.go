package mysql

import (
	"database/sql"
	"GoTest/pkg/models"
	"errors"
)

//добавление нового домена в базу данных
func (plm *PlantlistModel) InsertDomain(domain string) (int64, error) {
	result, err := plm.DB.Exec(`INSERT IGNORE INTO domain (id, name) VALUES (NULL, ?);`, domain)
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

//получение данных домена по его названию
func (plm *PlantlistModel) GetDomain(domName string) (*models.Domain, error) {
	selectDomain := `SELECT id, name FROM domain WHERE name = ?;`
	row := plm.DB.QueryRow(selectDomain, domName)

	dom := &models.Domain{}
	err := row.Scan(&dom.ID, &dom.Name)
	if err != nil {
		//если модели в программе и БД не совпадают
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}
	
	return dom, nil
}

//обновление данных по домену таксона
func (plm *PlantlistModel) UpdateDomain(domName string) error {
	_, err := plm.DB.Exec(`UPDATE domain SET name=? WHERE name=?;`, domName, domName)
	if err != nil {
		return err
	}

	return nil
}