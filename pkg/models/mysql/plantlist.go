package mysql

import (
	"database/sql"
	"GoTest/pkg/models"
	"errors"
)

//"time"

//обёртка пула подключений
type PlantlistModel struct {
	DB *sql.DB
}

type Systematic struct {
	Domain 		string 
	Kingdom 	string
	Department 	string
	Class 		string 
	Order 		string
	Family 		string 
	Genus 		string
}

type GPS struct {
	Latitude 	float64
	Longitude	float64
}

type SaveMeasure struct {
	Name    	string
	Description string
	//TODO: подумать
	Start		string
	End			string
}

var ErrNoRecord = errors.New("подходящей записи не найдено")

//сохранение растения в базу данных (используются плейсхолдеры для защиты от SQL-инъекций)
func (plm *PlantlistModel) InsertPlant(sys Systematic, status string, description string,
	publications []string, places []GPS, saveMeasures []SaveMeasure, rusName string, latinName string) error {

		//добавление новой систематики
		if _, err := plm.InsertClass(sys.Class); err != nil {
			return err
		}
		if _, err := plm.InsertDomain(sys.Domain); err != nil {
			return err
		}
		if _, err := plm.InsertKingdom(sys.Kingdom); err != nil {
			return err
		}
		if _, err := plm.InsertDepartment(sys.Department); err != nil {
			return err
		}
		if _, err := plm.InsertOrder(sys.Order); err != nil {
			return err
		}
		if _, err := plm.InsertFamily(sys.Family); err != nil {
			return err
		}
		if _, err := plm.InsertGenus(sys.Genus); err != nil {
			return err
		}
		if err := plm.InsertPublications(publications); err != nil {
			return err
		}
		if err := plm.InsertGPS(places); err != nil {
			return err
		}
		if err := plm.InsertConservations(saveMeasures); err != nil {
			return err
		}
		
		//получение систематики
		class, err := plm.GetClass(sys.Class)
		if err != nil {
			return err
		}
		domain, err := plm.GetDomain(sys.Domain)
		if err != nil {
			return err
		}
		kingdom, err := plm.GetKingdom(sys.Kingdom)
		if err != nil {
			return err
		}
		department, err := plm.GetDepartment(sys.Department)
		if err != nil {
			return err
		}
		order, err := plm.GetOrder(sys.Order)
		if err != nil {
			return err
		}
		family, err := plm.GetFamily(sys.Family)
		if err != nil {
			return err
		}
		genus, err := plm.GetGenus(sys.Genus)
		if err != nil {
			return err
		}
		stat, err := plm.GetStatus(status)
		if err != nil {
			return err
		}
		

		//добавление нового растения
		insertPlant := `INSERT INTO plant (id, russian_name, latin_name, description, id_class, id_department, id_domain, id_family, id_genus, id_kingdom, id_orger, id_status)
						VALUES (NULL, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?) ON DUPLICATE KEY UPDATE id=id;`
		_, err = plm.DB.Exec(insertPlant, rusName, latinName, description, class.ID, department.ID, domain.ID, family.ID, genus.ID, kingdom.ID, order.ID, stat.ID)
		if err != nil {
			return err
		}

		//получение растений
		plant, err := plm.GetPlant(rusName)
		if err != nil {
			return err
		}

		//получение точек из базы данных
		var points []*models.GPSPoint
		for _, point := range places {
			pointFromDB, err := plm.GetGPS(point.Longitude, point.Latitude)
			if err != nil {
				return err
			}
			points = append(points, pointFromDB)
		}

		//привязка растения к географическим точкам
		linkPlantAndGPS := "INSERT IGNORE INTO `plan-gps` (id, id_plant, id_gps) VALUES (NULL, ?, ?)"
		for _, point := range points {
			_, err = plm.DB.Exec(linkPlantAndGPS, plant.ID, point.ID)
		}

		//получение публикаций из базы данных
		var pubs []*models.Publication
		for _, pub := range publications {
			pubFromDB, err := plm.GetPublication(pub)
			if err != nil {
				return err
			}
			pubs = append(pubs, pubFromDB)
		}

		//привязка растения к публикациям
		linkPlantAndPub := "INSERT IGNORE INTO `plant-publlication` (id, id_plant, id_publlication) VALUES (NULL, ?, ?)"
		for _, pub := range pubs {
		_, err = plm.DB.Exec(linkPlantAndPub, plant.ID, pub.ID)
		}

		//получение мер сохранения из базы данных
		var conservations []*models.Conservation
		for _, cons := range saveMeasures {
			consFromDB, err := plm.GetConservation(cons.Name)
			if err != nil {
				return nil
			}
			conservations = append(conservations, consFromDB)
		}
		
		//привязка растения к мерам сохранения
		linkPlantAndConservatoin := "INSERT IGNORE INTO `plant-conservatioin` (id, id_plant, id_conservation) VALUES (NULL, ?, ?)"
		for _, cons := range conservations {
		_, err = plm.DB.Exec(linkPlantAndConservatoin, plant.ID, cons.ID)
		}
		
		return nil
}

//добавление нового класса в базу данных
func (plm *PlantlistModel) InsertClass(class string) (int64, error) {
	result, err := plm.DB.Exec(`INSERT IGNORE INTO class (id, name) VALUES (NULL, ?);`, class)
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

//получение данных класса по его названию
func (plm *PlantlistModel) GetClass(className string) (*models.Class, error) {
	selectClass := `SELECT id, name FROM class WHERE name = ?;`
	row := plm.DB.QueryRow(selectClass, className)

	class := &models.Class{}
	err := row.Scan(&class.ID, &class.Name)
	if err != nil {
		//если модели в программе и БД не совпадают
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}
	
	return class, nil
}

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

func (plm *PlantlistModel) InsertConservations(saveMeasures []SaveMeasure) error {
	insertConservation := `INSERT INTO conservation (id, name, description, start_date, end_date) VALUES (NULL, ?, ?, ?, ?) ON DUPLICATE KEY UPDATE id=id`
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

// func (PlantlistModel) GetAllPlants() {

// }

//возвращает растение по его рускоязычному названию
func (plm *PlantlistModel) GetPlant(rusName string) (*models.Plant, error) {
	selectPlant := `SELECT id, russian_name, latin_name, description, id_class, id_department, id_domain, id_family, id_genus, id_kingdom, id_orger, id_status
					 FROM plant WHERE russian_name = ?;`
	row := plm.DB.QueryRow(selectPlant, rusName)

	plant := &models.Plant{}
	err := row.Scan(&plant.ID, &plant.RusName, &plant.LatinName,
					&plant.Description, &plant.IDClass, &plant.IDDepartment, 
					&plant.IDDomain, &plant.IDFamily, &plant.IDGenus, 
					&plant.IDKingdom, &plant.IDOrder, &plant.IDStatus)
	if err != nil {
		//если модели в программе и БД не совпадают
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}
	
	return plant, nil
}