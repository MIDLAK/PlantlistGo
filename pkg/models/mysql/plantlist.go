package mysql

import (
	"database/sql"
)

//"time"
//"GoTest/pkg/models"

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

//сохранение растения в базу данных
func (plm *PlantlistModel) InsertPlant(sys Systematic, status string, description string,
	publications []string, places []GPS, saveMeasures []SaveMeasure) error {

		//используются плейсхолдеры для защиты от SQL-инъекций
		if err := plm.InsertClass(sys.Class); err != nil {
			return err
		}
		if err := plm.InsertDomain(sys.Domain); err != nil {
			return err
		}
		if err := plm.InsertKingdom(sys.Kingdom); err != nil {
			return err
		}
		if err := plm.InsertDepartment(sys.Department); err != nil {
			return err
		}
		if err := plm.InsertOrder(sys.Order); err != nil {
			return err
		}
		if err := plm.InsertFamily(sys.Family); err != nil {
			return err
		}
		if err := plm.InsertGenus(sys.Genus); err != nil {
			return err
		}
		if err := plm.InsertGenus(sys.Genus); err != nil {
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

		

		return nil
}

func (plm *PlantlistModel) InsertClass(class string) error {
	if _, err := plm.DB.Exec(`INSERT IGNORE INTO class (id, name) VALUES (NULL, ?);`, class); err != nil {
		return err
	}
	return nil
}

func (plm *PlantlistModel) InsertDepartment(department string) error {
	if _, err := plm.DB.Exec(`INSERT IGNORE INTO department (id, name) VALUES (NULL, ?);`, department); err != nil {
		return err
	}
	return nil
}

func (plm *PlantlistModel) InsertDomain(domain string) error {
	if _, err := plm.DB.Exec(`INSERT IGNORE INTO domain (id, name) VALUES (NULL, ?);`, domain); err != nil {
		return err
	}
	return nil
}

func (plm *PlantlistModel) InsertFamily(family string) error {
	if _, err := plm.DB.Exec(`INSERT IGNORE INTO family (id, name) VALUES (NULL, ?);`, family); err != nil {
		return err
	}
	return nil
}

func (plm *PlantlistModel) InsertGenus(genus string) error {
	if _, err := plm.DB.Exec(`INSERT IGNORE INTO genus (id, name) VALUES (NULL, ?);`, genus); err != nil {
		return err
	}
	return nil
}

func (plm *PlantlistModel) InsertKingdom(kingdom string) error {
	if _, err := plm.DB.Exec(`INSERT IGNORE INTO kingdom (id, name) VALUES (NULL, ?);`, kingdom); err != nil {
		return err
	}
	return nil
}

func (plm *PlantlistModel) InsertOrder(order string) error {
	//order - ключевое слово...
	if _, err := plm.DB.Exec("INSERT IGNORE INTO `order` (id, name) VALUES (NULL, ?);", order); err != nil {
		return err
	}
	return nil
}

func (plm *PlantlistModel) InsertGPS(places []GPS) error {
	insertGPS := `INSERT IGNORE INTO gps (latitude, longitude) VALUES (?, ?);`
	for _, point := range places {
		if _, err := plm.DB.Exec(insertGPS, point.Latitude, point.Longitude); err != nil {
			return err
		}
	}
	return nil
}

func (plm *PlantlistModel) InsertPublications(publications []string) error {
	insertPublication := `INSERT IGNORE INTO publication (bibliographic_description) VALUES (?);`
	for _, publication := range publications {
		if _, err := plm.DB.Exec(insertPublication, publication); err != nil {
			return err
		}
	}
	return nil
}

func (plm *PlantlistModel) InsertConservations(saveMeasures []SaveMeasure) error {
	insertConservation := `INSERT INTO conservation (name, description, start_date, end_date) VALUES (?, ?, ?, ?) ON DUPLICATE KEY UPDATE id=id`
	for _, measure := range saveMeasures {
		if _, err := plm.DB.Exec(insertConservation, measure.Name, measure.Start, measure.End); err != nil {
			return err
		} 
	}
	return nil
}
// func (PlantlistModel) GetAllPlants() {

// }

// func (PlantLlistModel) GetPlant() {

// }