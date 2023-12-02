package mysql

import (
	"database/sql"
	"GoTest/pkg/models"
	"errors"
)

type GPS struct {
	Latitude 	float64
	Longitude	float64
}

func (plm *PlantlistModel) InsertGPS(places []GPS) error {
	insertGPS := `INSERT IGNORE INTO gps (id, latitude, longitude) VALUES (NULL, ?, ?);`
	for _, point := range places {
		if _, err := plm.DB.Exec(insertGPS, point.Latitude, point.Longitude); err != nil {
			return err
		}
	}
	return nil
}

//получение GPS-координаты по широте и долготе
func (plm *PlantlistModel) GetGPS(longitude float64, latitude float64) (*models.GPSPoint, error) {
	selectGPS := `SELECT id, latitude, longitude 
				  FROM gps WHERE ABS(latitude - ?) < 0.001 AND ABS(longitude - ?) < 0.001;`
	row := plm.DB.QueryRow(selectGPS, latitude, longitude)

	point := &models.GPSPoint{}
	err := row.Scan(&point.ID, &point.Latitude, &point.Longitude)
	if err != nil {
		//если модели в программе и БД не совпадают
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}
	
	return point, nil
}