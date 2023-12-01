package main

import (
	"html/template"
	"net/http"
	"encoding/json"
	"GoTest/pkg/models/mysql"
	"strconv"
	"errors"
	"GoTest/pkg/models"
)

type Previews struct {
	List	[]*models.PlantPreview
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

	//превью всех таксонов
	previews, err := app.dbPlantlist.GetPlantsPreview()
	if err != nil {
		app.serverError(w, err)
	}

	//запись шаблона в тело HTTP запроса
	err = ts.Execute(w, Previews{List: previews})
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

//обработчки /about
func (app *application) about(w http.ResponseWriter, r *http.Request) {
	files := []string{
		"./ui/html/about.page.html",
		"./ui/html/base.layout.html",
		"./ui/html/footer.partial.html",
		"./ui/html/place.partial.html",
	}

	//чтение файла из шаблона
	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.serverError(w, err)
		return
	}

	//получаем данные, что за растение
	r.ParseForm()
	plantName := r.Form.Get("plant-name-add")
	app.infoLog.Printf("Запрос данных по таксону <<%v>>", plantName)

	//формирование тела шаблона
	aboutPlant, err := app.dbPlantlist.GetAboutPlant(plantName)
	if err != nil {
		app.serverError(w, err)
		return
	}

	//запись шаблона в тело HTTP запроса
	err = ts.Execute(w, aboutPlant)
	if err != nil {
		app.serverError(w, err)
		return
	}
}

//обработчик /model
func (app *application) model(w http.ResponseWriter, r *http.Request) {
	files := []string{
		"./ui/html/add-or-edit.page.html",
		"./ui/html/base.layout.html",
		"./ui/html/footer.partial.html",
		"./ui/html/place.partial.html",
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

//структура, соответствующая JSON для обмена
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

//обработчик /newplant (только для POST-запроса)
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

type Place struct {
	Latitude	float64
	Longitude	float64
}

type SM struct {
	SaveName	string
	Description	string
	Start   	string
	End     	string
}

//структура, соответствующая JSON для обмена
type PlantDataOutput struct {
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
	Places       	[]Place  `json:"places"`
	SaveMeasure 	[]SM	 `json:"saveMeasure"`
}

type PlantNameInput struct {
	PlantName		string
}

//обработчик /getplant (только для POST-запроса)
func (app *application) getplant(w http.ResponseWriter, r *http.Request) {
	//если кто-то просто пытается перейти по /newpalnt
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}

	//получение данных
	var plantNameInput PlantNameInput
	err := json.NewDecoder(r.Body).Decode(&plantNameInput)
	if err != nil {
		app.serverError(w, err)
	}

	//формироване ответа
	var outputData PlantDataOutput
	aboutPlant, err := app.dbPlantlist.GetAboutPlant(plantNameInput.PlantName)
	if err != nil {
		outputData = PlantDataOutput{Name: "", LatinName: "",
			Domain: "", Kingdom: "",
			Department: "", Class: "",
			Order: "", Family: "",
			Genus: "", Status: "",
			Description: "", Publications: []string{""}}
	} else {
		outputData = PlantDataOutput{Name: aboutPlant.Name, LatinName: aboutPlant.LatinName,
									Domain: aboutPlant.Domain, Kingdom: aboutPlant.Kingdom,
									Department: aboutPlant.Department, Class: aboutPlant.Class,
									Order: aboutPlant.Order, Family: aboutPlant.Family,
									Genus: aboutPlant.Genus, Status: aboutPlant.Status,
									Description: aboutPlant.Description, Publications: aboutPlant.Publications}
		for _, point := range aboutPlant.Places {
			p := Place{Latitude: point.Latitude, Longitude: point.Longitude}
			outputData.Places = append(outputData.Places, p)
		}
		for _, conservation := range aboutPlant.SaveMeasures {
			sm := SM{SaveName: conservation.Name, Description: conservation.Description, Start: conservation.Start.Format("2006-09-02"),
							End: conservation.End.Format("2006-09-02")} 
			outputData.SaveMeasure = append(outputData.SaveMeasure, sm)
		}
	}
	

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(outputData)

}

//обработчик /edit
func (app *application) editplant(w http.ResponseWriter, r *http.Request) {
	files := []string{
		"./ui/html/add-or-edit.page.html",
		"./ui/html/base.layout.html",
		"./ui/html/footer.partial.html",
		"./ui/html/place.partial.html",
	}

	//чтение файла из шаблона
	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.serverError(w, err)
		return
	}

	//получаем данные, что за растение
	r.ParseForm()
	plantName := r.Form.Get("plant-name-edit")
	app.infoLog.Printf("Запрос данных по таксону <<%v>>", plantName)

	//формирование тела шаблона
	aboutPlant, err := app.dbPlantlist.GetAboutPlant(plantName)
	if err != nil {
		app.serverError(w, err)
		return
	}

	//запись шаблона в тело HTTP запроса
	err = ts.Execute(w, aboutPlant)
	if err != nil {
		app.serverError(w, err)
		return
	}
}

//обработчик /updateplant (только для POST-запроса)
// func (app *application) updateplant(w http.ResponseWriter, r *http.Request) {
// 	//добавление (обновление) таксона
// 	//если кто-то просто пытается перейти по /newpalnt
// 	if r.Method != http.MethodPost {
// 		w.Header().Set("Allow", http.MethodPost)
// 		app.clientError(w, http.StatusMethodNotAllowed)
// 		return
// 	}

// 	//получение данных
// 	var plantDataInput PlantDataInput
// 	err := json.NewDecoder(r.Body).Decode(&plantDataInput)
// 	if err != nil {
// 		app.serverError(w, err)
// 	}

// 	//расшифровка полученных данных
// 	sys :=  mysql.Systematic {Domain: plantDataInput.Domain, Kingdom: plantDataInput.Kingdom,
// 							  Class: plantDataInput.Class, Order: plantDataInput.Order, 
// 							  Family: plantDataInput.Family, Genus: plantDataInput.Genus,
// 							  Department: plantDataInput.Department}
// 	var places []mysql.GPS
// 	for _, elem := range plantDataInput.Places {
// 		latitude, err := strconv.ParseFloat(elem.Latitude, 64);
// 		if err != nil {
// 			app.serverError(w, err)
// 		}
// 		longitude, err := strconv.ParseFloat(elem.Longitude, 64);
// 		if err != nil {
// 			app.serverError(w, err)
// 		}
// 		places = append(places, mysql.GPS{Latitude: latitude, Longitude: longitude})
// 	}
// 	var saveMeasures []mysql.SaveMeasure
// 	for _, elem := range plantDataInput.SaveMeasure {
// 		saveMeasures = append(saveMeasures, mysql.SaveMeasure{Name: elem.SaveName, Description: elem.Description, 
// 						 Start: elem.StartDate, End: elem.EndDate})
// 	}

// 	//добавление (обновление) таксона
// 	err = app.dbPlantlist.UpdatePlant(sys, plantDataInput.Status, plantDataInput.Description,
// 									   plantDataInput.Publications, places, saveMeasures,
// 									   plantDataInput.Name, plantDataInput.LatinName)
// 	if err != nil {
// 		if errors.Is(err, mysql.ErrNoRecord) {
// 			app.notFound(w)
// 		} else {
// 			app.serverError(w, err)
// 		}
// 	} else {
// 		app.infoLog.Printf("Обновлена информация о таксоне <<%v>>", plantDataInput.Name)
// 	}
// }