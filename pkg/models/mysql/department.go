package mysql

import (
	"database/sql"
	"GoTest/pkg/models"
	"errors"
)

//добавление нового отдела в базу данных
func (plm *PlantlistModel) InsertDepartment(department string) (int64, error) {
	result, err := plm.DB.Exec(`INSERT IGNORE INTO department (id, name) VALUES (NULL, ?);`, department)
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

//получение данных отдела по его названию
func (plm *PlantlistModel) GetDepartment(depName string) (*models.Department, error) {
	selectDepartment := `SELECT id, name FROM department WHERE name = ?;`
	row := plm.DB.QueryRow(selectDepartment, depName)

	dep := &models.Department{}
	err := row.Scan(&dep.ID, &dep.Name)
	if err != nil {
		//если модели в программе и БД не совпадают
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}
	
	return dep, nil
}

//обновление данных по отделу таксона
func (plm *PlantlistModel) UpdateDepartment(depName string) error {
	_, err := plm.DB.Exec(`UPDATE department SET name=? WHERE name=?;`, depName, depName)
	if err != nil {
		return err
	}

	return nil
}