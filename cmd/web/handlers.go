package main

import (
	"html/template"
	"net/http"
	"encoding/json"
	"GoTest/pkg/models/mysql"
	"strconv"
	"errors"
)

// "GoTest/pkg/models"
// 


type Card struct {
	PlantName	string
	PlantText	string
}

type Cards struct {
	CardsList	[]Card
}

//обработчик /
func (app *application) home(w http.ResponseWriter, r *http.Request) {
	//реакция на несуществующие страницы
	if r.URL.Path != "/" {
		app.notFound(w)
		return
	}

	// Инициализируем срез содержащий пути к файлам
	// файл home.page.tmpl должен быть *первым* файлом в срезе
	files := []string{
		"./ui/html/home.page.html",
		"./ui/html/base.layout.html",
		"./ui/html/footer.partial.html",
        "./ui/html/card.partial.html",
	}

	//чтение файла из шаблона
	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.serverError(w, err)
		return
	}

	//тестовые данные
	cladoniaCard := Card{PlantName: "Кладония", PlantText: "Лишайник такой"}
	smolCard := Card{PlantName: "Смолёвка", PlantText: "Лишайник другой"}
	cards := Cards{}
	cards.CardsList = append(cards.CardsList, cladoniaCard, smolCard)

	//запись шаблона в тело HTTP запроса
	err = ts.Execute(w, cards)
	if err != nil {
		app.serverError(w, err)
		return
	}
}

//обработчик /details
func (app *application) details(w http.ResponseWriter, r *http.Request) {
	files := []string{
		"./ui/html/details.page.html",
		"./ui/html/base.layout.html",
		"./ui/html/footer.partial.html",
	}

	//чтение файла из шаблона
	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.serverError(w, err)
		return
	}

	//запись шаблона в тело HTTP запроса
	err = ts.Execute(w, nil)
	if err != nil {
		app.serverError(w, err)
		return
	}
}

//обработчик /model
func (app *application) model(w http.ResponseWriter, r *http.Request) {
	files := []string{
		"./ui/html/edit-and-add.page.html",
		"./ui/html/base.layout.html",
		"./ui/html/footer.partial.html",
	}

	//чтение файла из шаблона
	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.serverError(w, err)
		return
	}

	//запись шаблона в тело HTTP запроса
	err = ts.Execute(w, nil)
	if err != nil {
		app.serverError(w, err)
		return
	}
}

//структура, соответствующая принимаемому JSON
type PlantDataInput struct {
	Name         	string   `json:"name"`
	LatinName		string	 `json:"latinName"`
	Domain       	string   `json:"domain"`
	Kingdom      	string   `json:"kingdom"`
	Department   	string   `json:"department"`
	Class        	string   `json:"class"`
	Order        	string   `json:"order"`
	Family       	string   `json:"family"`
	Genus        	string   `json:"genus"`
	Status       	string   `json:"status"`
	Description  	string   `json:"description"`
	Publications 	[]string `json:"publications"`
	Places       	[]struct {
		Latitude  string `json:"latitude"`
		Longitude string `json:"longitude"`
	} 						 `json:"places"`
	SaveMeasure 	[]struct {
		SaveName    string `json:"saveName"`
		Description string `json:"description"`
		StartDate   string `json:"startDate"`
		EndDate     string `json:"endDate"`
	} 						 `json:"saveMeasure"`
}

//обработчик /newplant
func (app *application) newplant(w http.ResponseWriter, r *http.Request) {
	//если кто-то просто пытается перейти по /newpalnt
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}

	//получение данных
	var plantDataInput PlantDataInput
	err := json.NewDecoder(r.Body).Decode(&plantDataInput)
	if err != nil {
		app.serverError(w, err)
	}

	//расшифровка полученных данных
	sys :=  mysql.Systematic {Domain: plantDataInput.Domain, Kingdom: plantDataInput.Kingdom,
							  Class: plantDataInput.Class, Order: plantDataInput.Order, 
							  Family: plantDataInput.Family, Genus: plantDataInput.Genus,
							  Department: plantDataInput.Department}
	var places []mysql.GPS
	for _, elem := range plantDataInput.Places {
		latitude, err := strconv.ParseFloat(elem.Latitude, 64);
		if err != nil {
			app.serverError(w, err)
		}
		longitude, err := strconv.ParseFloat(elem.Longitude, 64);
		if err != nil {
			app.serverError(w, err)
		}
		places = append(places, mysql.GPS{Latitude: latitude, Longitude: longitude})
	}
	var saveMeasures []mysql.SaveMeasure
	for _, elem := range plantDataInput.SaveMeasure {
		saveMeasures = append(saveMeasures, mysql.SaveMeasure{Name: elem.SaveName, Description: elem.Description, 
						 Start: elem.StartDate, End: elem.EndDate})
	}

	//добавление (обновление) таксона
	err = app.dbPlantlist.InsertPlant(sys, plantDataInput.Status, plantDataInput.Description,
									   plantDataInput.Publications, places, saveMeasures,
									   plantDataInput.Name, plantDataInput.LatinName)
	if err != nil {
		if errors.Is(err, mysql.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
	} else {
		app.infoLog.Printf("Добавлена / обновлена информация о таксоне <<%v>>", plantDataInput.Name)
	}
}
// type WeatherDataInput struct {
// 	Speed                    string
// 	Direction                string
// 	AirHumidity              string
// 	Temperature              string
// 	AtmosphericPrecipitation string
// }

// type OutputData struct {
// 	Wind                     models.Wind
// 	Temperature              models.Temperature
// 	AirHumidity              models.AirHumidity
// 	AtmosphericPrecipitation models.AtmosphericPrecipitation
// 	TimeInterval             int
// 	FireCoefficient          int
// 	FireClass                int
// }

// func (app *application) getWeather(w http.ResponseWriter, r *http.Request) {
// 	if r.Method != http.MethodPost {
// 		w.Header().Set("Allow", http.MethodPost)
// 		app.clientError(w, http.StatusMethodNotAllowed)
// 		return
// 	}

// 	//получение данных
// 	var weatherInput WeatherDataInput
// 	err := json.NewDecoder(r.Body).Decode(&weatherInput)
// 	if err != nil {
// 		app.serverError(w, err)
// 	}

// 	/*Получение вводных*/
// 	//ветер
// 	windOutput := models.NewWind()
// 	speed, err := strconv.Atoi(weatherInput.Speed)
// 	if err != nil {
// 		app.serverError(w, err)
// 		return
// 	}
// 	windOutput.Speed = speed

// 	switch weatherInput.Direction {
// 	case "North":
// 		windOutput.Direction = models.North
// 	case "East":
// 		windOutput.Direction = models.East
// 	case "South":
// 		windOutput.Direction = models.South
// 	case "West":
// 		windOutput.Direction = models.West
// 	}

// 	//температура
// 	temperature, err := strconv.Atoi(weatherInput.Temperature)
// 	if err != nil {
// 		app.serverError(w, err)
// 		return
// 	}
// 	temperatureOutput := models.Temperature{
// 		DegreesCelsius: temperature,
// 	}

// 	//влажность
// 	humidity, err := strconv.Atoi(weatherInput.AirHumidity)
// 	if err != nil {
// 		app.serverError(w, err)
// 		return
// 	}
// 	airHumidityOutput := models.AirHumidity{
// 		Percent: humidity,
// 	}

// 	//осадки
// 	precipitation, err := strconv.Atoi(weatherInput.AtmosphericPrecipitation)
// 	if err != nil {
// 		app.serverError(w, err)
// 		return
// 	}
// 	atmosphericPrecipitationOutput := models.AtmosphericPrecipitation{
// 		Millimeters: precipitation,
// 	}

// 	//погода в целом
// 	weatherOutput := models.Weather{
// 		Wind:                     windOutput,
// 		Temperature:              temperatureOutput,
// 		AirHumidity:              airHumidityOutput,
// 		AtmosphericPrecipitation: atmosphericPrecipitationOutput,
// 	}

// 	/* Формирование ответа браузерной части приложения */
// 	cof, class := weatherOutput.FireDangerCoefficientAndClass()
// 	q, _ := atmosphericPrecipitationOutput.OverwhelmingPrecipitation(cof)
// 	interval := 3000 - windOutput.Speed*100 + q*100
// 	outputData := OutputData{
// 		Wind:                     windOutput,
// 		Temperature:              temperatureOutput,
// 		AirHumidity:              airHumidityOutput,
// 		AtmosphericPrecipitation: atmosphericPrecipitationOutput,
// 		TimeInterval:             interval,
// 		FireCoefficient:          cof,
// 		FireClass:                class,
// 	}

// 	w.Header().Set("Content-Type", "application/json")
// 	_ = json.NewEncoder(w).Encode(outputData)
// }
