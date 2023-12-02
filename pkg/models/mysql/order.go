package mysql

import (
	"database/sql"
	"GoTest/pkg/models"
	"errors"
)

//добавление нового порядка в базу данных
func (plm *PlantlistModel) InsertOrder(order string) (int64, error) {
	//order - ключевое слово...
	result, err := plm.DB.Exec("INSERT IGNORE INTO `order` (id, name) VALUES (NULL, ?);", order)
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

//получение данных порядка по его названию
func (plm *PlantlistModel) GetOrder(orderName string) (*models.Order, error) {
	selectOrder := "SELECT id, name FROM `order` WHERE name = ?;"
	row := plm.DB.QueryRow(selectOrder, orderName)

	order := &models.Order{}
	err := row.Scan(&order.ID, &order.Name)
	if err != nil {
		//если модели в программе и БД не совпадают
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}
	
	return order, nil
}

//обновление данных по семейству таксона
func (plm *PlantlistModel) UpdateOrder(orderName string) error {
	_, err := plm.DB.Exec("UPDATE `order` SET name=? WHERE name=?;`, orderName, orderName")
	if err != nil {
		return err
	}

	return nil
}