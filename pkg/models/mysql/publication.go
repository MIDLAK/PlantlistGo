package mysql

import (
	"database/sql"
	"GoTest/pkg/models"
	"errors"
)

//добавление списка публикаций в БД
func (plm *PlantlistModel) InsertPublications(publications []string) error {
	insertPublication := `INSERT IGNORE INTO publlication (id, link, bibliographic_description, description) VALUES (NULL, NULL, ?, NULL);`
	for _, publication := range publications {
		if _, err := plm.DB.Exec(insertPublication, publication); err != nil {
			return err
		}
	}
	return nil
}

//получение публикации из БД по содержанию
func (plm *PlantlistModel) GetPublication(bibilio string) (*models.Publication, error) {
	selectPublication := `SELECT id, link, bibliographic_description, description
						  FROM publlication WHERE bibliographic_description = ?;`
	row := plm.DB.QueryRow(selectPublication, bibilio)

	pub := &models.Publication{}
	err := row.Scan(&pub.ID, &pub.Link, &pub.Biblio, &pub.Description)
	if err != nil {
		//если модели в программе и БД не совпадают
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}
	
	return pub, nil						  
}