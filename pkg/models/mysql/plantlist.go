package mysql

import (
	"database/sql"
	"GoTest/pkg/models"
	"errors"
	"time"
)

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

var ErrNoRecord = errors.New("Подходящей записи не найдено")

//сохранение растения в базу данных (используются плейсхолдеры для защиты от SQL-инъекций)
func (plm *PlantlistModel) InsertPlant(sys Systematic, status string, description string,
	publications []string, places []GPS, saveMeasures []SaveMeasure, rusName string, latinName string) error {


		plantFromDB, err := plm.GetAboutPlant(rusName)
		//если такое растение уже существует, то разрыв связей
		//с геогр. координатами, публикациями и мерами сохранения
		if err == nil {
			plm.deletePlant(plantFromDB.ID)
		}

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


//изменение данных по растению в базе данных
// func (plm *PlantlistModel) UpdatePlant(sys Systematic, status string, description string,
// 	publications []string, places []GPS, saveMeasures []SaveMeasure, rusName string, latinName string) error {
	
// 	//получение данных по растению
// 	plant, err := plm.GetAboutPlant(rusName)
// 	if err != nil {
// 		return err
// 	}

// 	//разрыв связей с известными параметрами
// 	err = plm.deletePlant(rusName)
// 	if err != nil {
// 		return err
// 	}



// 	return nil
// }

//удаление растения
func (plm *PlantlistModel) deletePlant(plantID int) error {
	//разрыв связей с координатам
	_, err := plm.DB.Exec("DELETE FROM `plan-gps` WHERE id_plant=?", plantID)
	if err != nil {
		return err
	}

	//разрыв связей с публикациями
	_, err = plm.DB.Exec("DELETE FROM `plant-publlication` WHERE id_plant=?", plantID)
	if err != nil {
		return err
	}

	//разрыв связей с мерами сохранения
	_, err = plm.DB.Exec("DELETE FROM `plant-conservatioin` WHERE id_plant=?", plantID)
	if err != nil {
		return err
	}

	//удаление растения
	_, err = plm.DB.Exec("DELETE FROM plant WHERE id=?", plantID)
	if err != nil {
		return err
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

//обновление класса для растения
func (plm *PlantlistModel) UpdateClass(className string) error {
	_, err := plm.DB.Exec(`UPDATE class SET name=? WHERE name=?;`, className)
	if err != nil {
		return err
	}

	return nil
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

//обновление данных по отделу таксона
func (plm *PlantlistModel) UpdateDepartment(depName string) error {
	_, err := plm.DB.Exec(`UPDATE department SET name=? WHERE name=?;`, depName, depName)
	if err != nil {
		return err
	}

	return nil
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

//обновление данных по домену таксона
func (plm *PlantlistModel) UpdateDomain(domName string) error {
	_, err := plm.DB.Exec(`UPDATE domain SET name=? WHERE name=?;`, domName, domName)
	if err != nil {
		return err
	}

	return nil
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

//обновление данных по семейству таксона
func (plm *PlantlistModel) UpdateFamily(famName string) error {
	_, err := plm.DB.Exec(`UPDATE family SET name=? WHERE name=?;`, famName, famName)
	if err != nil {
		return err
	}

	return nil
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

//обновление данных по семейству таксона
func (plm *PlantlistModel) UpdateGenus(genName string) error {
	_, err := plm.DB.Exec(`UPDATE genus SET name=? WHERE name=?;`, genName, genName)
	if err != nil {
		return err
	}

	return nil
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

//обновление данных по семейству таксона
func (plm *PlantlistModel) UpdateKingdom(kinName string) error {
	_, err := plm.DB.Exec(`UPDATE kingdom SET name=? WHERE name=?;`, kinName, kinName)
	if err != nil {
		return err
	}

	return nil
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

//обновление данных по семейству таксона
func (plm *PlantlistModel) UpdateOrder(orderName string) error {
	_, err := plm.DB.Exec("UPDATE `order` SET name=? WHERE name=?;`, orderName, orderName")
	if err != nil {
		return err
	}

	return nil
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

//возвращает только название и описание для растений
func (plm *PlantlistModel) GetPlantsPreview() ([]*models.PlantPreview, error) {
	selectPreview := `SELECT plant.russian_name, plant.description FROM plant;`
	rows, err := plm.DB.Query(selectPreview)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	//просмотр полученных данных
	var previews []*models.PlantPreview
	for rows.Next() {
		prv := &models.PlantPreview{}
		err = rows.Scan(&prv.Name, &prv.Description)
		if err != nil {
			return nil, err
		}

		previews = append(previews, prv)
	}

	//дополнительная проверка
	if err = rows.Err(); err != nil {
	 	return nil, err
	}

	return previews, nil
}

//как это будет получено из базы
type AboutPlantRow struct {
	ID				int
	RusName			string
	LatName			string
	Description		string
	Class			string
	Department		string
	Domain			string
	Family			string
	Genus			string
	Kingdom			string
	Order			string
	Status			string
	Latitude		sql.NullFloat64
	Longitude		sql.NullFloat64
	Conservation	sql.NullString
	ConsDesc		sql.NullString
	Start			sql.NullTime
	End				sql.NullTime
	Biblio			sql.NullString
}

//чем полученные из БД данные должны стать
type Place struct {
	Latitude	float64
	Longitude	float64
}

type SM struct {
	Name		string
	Description string
	Start   	time.Time
	End     	time.Time
}

type AboutPlant struct {
	ID				int
	Name         	string
	LatinName		string
	Domain       	string
	Kingdom      	string
	Department   	string
	Class        	string
	Order        	string
	Family       	string
	Genus        	string
	Status       	string
	Description  	string
	Publications 	[]string 
	Places       	[]Place					 
	SaveMeasures 	[]SM						 
}

//TODO: тут творится кошмар
func (plm *PlantlistModel) GetAboutPlant(rusname string) (*AboutPlant, error) {
	selectAboutPlants := "SELECT p.id, p.russian_name AS rus, p.latin_name, p.description, cls.name AS class, dep.name AS department, dom.name AS domain, f.name AS family, g.name AS genus, k.name AS kingdom, ord.name AS `order`, s.name AS `status`,  gps.latitude, gps.longitude, cons.name AS conservation, cons.description, cons.start_date, cons.end_date, pub.bibliographic_description AS biblio FROM plant AS p LEFT JOIN class 	 AS cls	 ON p.id_class = cls.id LEFT JOIN department AS dep  ON p.id_department = dep.id LEFT JOIN domain 	 AS dom	 ON p.id_domain = dom.id LEFT JOIN family 	 AS f 	 ON p.id_family = f.id LEFT JOIN genus 	 AS g 	 ON p.id_genus = g.id LEFT JOIN kingdom 	 AS k 	 ON p.id_kingdom = k.id LEFT JOIN `order` 	 AS ord	 ON p.id_orger = ord.id LEFT JOIN `status`	 AS s	 ON p.id_status = s.id LEFT JOIN `plan-gps` AS pgps ON p.id = pgps.id_plant LEFT JOIN gps ON gps.id = pgps.id_gps LEFT JOIN `plant-conservatioin` AS pcons ON  p.id = pcons.id_plant LEFT JOIN conservation			AS cons  ON  cons.id = pcons.id_conservation LEFT JOIN `plant-publlication`  AS ppub	 ON	p.id = ppub.id_plant LEFT JOIN publlication			AS pub 	 ON pub.id = ppub.id_publlication;"
	
	rows, err := plm.DB.Query(selectAboutPlants)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	

	//просмотр полученных данных
	var plantRows []*AboutPlantRow
	for rows.Next() {
		plt := &AboutPlantRow{}
		err = rows.Scan(&plt.ID, &plt.RusName, &plt.LatName,
			&plt.Description, &plt.Class, &plt.Department,
			&plt.Domain, &plt.Family, &plt.Genus, &plt.Kingdom,
			&plt.Order, &plt.Status, &plt.Latitude, &plt.Longitude,
			&plt.Conservation, &plt.ConsDesc, &plt.Start, &plt.End, 
			&plt.Biblio)
		if err != nil {
			return nil, err
		}

		plantRows = append(plantRows, plt)
	}

	//преобразование данных
	prevID := 0
	var aboutPlants []*AboutPlant
	for _, r := range plantRows {
		//если новое растение
		if r.ID != prevID {
			prevID = r.ID
			//проверка, что описание такого растения ещё не начато
			for _, about := range aboutPlants {
				//если начато, то добавляем данные
				if about.ID == r.ID {
					if r.Biblio.Valid {
						flag := true
						for _, pub := range about.Publications {
							if pub == r.Biblio.String {
								flag = false
							}
						}

						if flag {
							about.Publications = append(about.Publications, r.Biblio.String)
						}
					}
					
					if r.Latitude.Valid {
						flag := true
						for _, place := range about.Places {
							if place.Latitude == r.Latitude.Float64 && place.Longitude == r.Longitude.Float64  {
								flag = false
							}
						}

						if flag {
							about.Places = append(about.Places, Place{Latitude: float64(r.Latitude.Float64), 
															Longitude: float64(r.Longitude.Float64)})
						}
					}

					if r.Conservation.Valid && r.Start.Valid && r.End.Valid {
						flag := true
						for _, conservation := range about.SaveMeasures {
							if conservation.Name == r.Conservation.String {
								flag = false
							}
						}

						if flag {
						about.SaveMeasures = append(about.SaveMeasures, SM{Name: r.Conservation.String,
																	Description: r.ConsDesc.String,
																	Start: r.Start.Time,
																	End:	r.End.Time})
						}
					}
				}
			}

			biblio := ""
			if r.Biblio.Valid {
				biblio = r.Biblio.String 
			}
			latitude := 0.0
			longitude := 0.0
			if r.Latitude.Valid && r.Longitude.Valid {
				latitude = r.Latitude.Float64
				longitude = r.Longitude.Float64
			}
			consName := ""
			if r.Conservation.Valid {
				consName = r.Conservation.String
			}
			description := ""
			if r.ConsDesc.Valid {
				description = r.ConsDesc.String
			}
			

			
			aboutPlants = append(aboutPlants, &AboutPlant{ID: r.ID, Name: r.RusName, LatinName: r.LatName,
								Domain: r.Domain, Kingdom: r.Kingdom, Department: r.Department,
								Class: r.Class, Order: r.Order, Family: r.Family, Genus: r.Genus,
								Status: r.Status, Description: r.Description, 
								Publications: []string{biblio}, 
								Places: []Place{{Latitude: latitude, Longitude: longitude}},
								SaveMeasures: []SM{{Name: consName, Description: description,
													End: r.End.Time, Start: r.Start.Time}}})
		} else {
			for _, about := range aboutPlants {
				//если начато, то добавляем данные
				if about.ID == r.ID {
					if r.Biblio.Valid {
						flag := true
						for _, pub := range about.Publications {
							if pub == r.Biblio.String {
								flag = false
							}
						}

						if flag {
							about.Publications = append(about.Publications, r.Biblio.String)
						}
					}
					
					if r.Latitude.Valid {
						flag := true
						for _, place := range about.Places {
							if place.Latitude == r.Latitude.Float64 && place.Longitude == r.Longitude.Float64  {
								flag = false
							}
						}

						if flag {
							about.Places = append(about.Places, Place{Latitude: float64(r.Latitude.Float64), 
															Longitude: float64(r.Longitude.Float64)})
						}
					}

					if r.Conservation.Valid && r.Start.Valid && r.End.Valid {
						flag := true
						for _, conservation := range about.SaveMeasures {
							if conservation.Name == r.Conservation.String {
								flag = false
							}
						}

						if flag {
						about.SaveMeasures = append(about.SaveMeasures, SM{Name: r.Conservation.String,
																	Description: r.ConsDesc.String,
																	Start: r.Start.Time,
																	End:	r.End.Time})
						}
					}
				}
			}
		}
	}

	//дополнительная проверка
	if err = rows.Err(); err != nil {
		return nil, err
	}

	for _, elem := range aboutPlants {
		if elem.Name == rusname {
			return elem, nil
		}
	}

	//если ничего не нашлось
	return nil, ErrNoRecord
}


//очень страшно
/*
SELECT p.id, p.russian_name rus, p.latin_name, p.description, 
cls.name AS class, dep.name AS department, dom.name AS domain, f.name AS family,
g.name AS genus, k.name AS kingdom, ord.name AS `order`, s.name AS `status`, 
gps.latitude, gps.longitude, cons.name AS conservation, cons.start_date, cons.end_date,
pub.bibliographic_description AS biblio

FROM plant AS p

-- Систематика
LEFT JOIN class 	 AS cls	 ON p.id_class = cls.id
LEFT JOIN department AS dep  ON p.id_department = dep.id
LEFT JOIN domain 	 AS dom	 ON p.id_domain = dom.id
LEFT JOIN family 	 AS f 	 ON p.id_family = f.id
LEFT JOIN genus 	 AS g 	 ON p.id_genus = g.id
LEFT JOIN kingdom 	 AS k 	 ON p.id_kingdom = k.id
LEFT JOIN `order` 	 AS ord	 ON p.id_orger = ord.id
LEFT JOIN `status`	 AS s	 ON p.id_status = s.id

-- GPS-координаты
LEFT JOIN `plan-gps` AS pgps ON p.id = pgps.id_plant
LEFT JOIN gps 				 ON gps.id = pgps.id_gps

-- Меры защиты
LEFT JOIN `plant-conservatioin` AS pcons ON  p.id = pcons.id_plant
LEFT JOIN conservation			AS cons  ON  cons.id = pcons.id_conservation

-- Публикации
LEFT JOIN `plant-publlication`  AS ppub	 ON	p.id = ppub.id_plant
LEFT JOIN publlication			AS pub 	 ON pub.id = ppub.id_publlication;	
*/
