package mysql

import (
	"database/sql"
	"GoTest/pkg/models"
	"errors"
)

//получение статуса растения
func (plm *PlantlistModel) GetStatus(statusName string) (*models.Status, error) {
	selectStatus := "SELECT id, name, description FROM `status` WHERE name = ?;"
	row := plm.DB.QueryRow(selectStatus, statusName)

	status := &models.Status{}
	err := row.Scan(&status.ID, &status.Name, &status.Description)
	if err != nil {
		//если модели в программе и БД не совпадают
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}
	
	return status, nil
}