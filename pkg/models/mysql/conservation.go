package mysql

import (
	"database/sql"
	"GoTest/pkg/models"
	"errors"
)

type SaveMeasure struct {
	Name    	string
	Description string
	Start		string
	End			string
}

func (plm *PlantlistModel) InsertConservations(saveMeasures []SaveMeasure) error {
	insertConservation := `INSERT IGNORE INTO conservation (id, name, description, start_date, end_date) VALUES (NULL, ?, ?, ?, ?)`
	for _, measure := range saveMeasures {
		if _, err := plm.DB.Exec(insertConservation, measure.Name, measure.Description, measure.Start, measure.End); err != nil {
			return err
		} 
	}
	return nil
}

//получение публикации из БД по содержанию
func (plm *PlantlistModel) GetConservation(consName string) (*models.Conservation, error) {
	selectConservation := `SELECT id, name, description, start_date, end_date
						  FROM conservation WHERE name = ?;`
	row := plm.DB.QueryRow(selectConservation, consName)

	cons := &models.Conservation{}
	err := row.Scan(&cons.ID, &cons.Name, &cons.Description, &cons.Start, &cons.End)
	if err != nil {
		//если модели в программе и БД не совпадают
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}
	
	return cons, nil						  
}