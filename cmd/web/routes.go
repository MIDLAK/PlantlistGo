package main

import "net/http"

//регистрация маршрутов
func (app *application) routes() *http.ServeMux {
	mux := http.NewServeMux() //создание мультиплексора
	//регистрация обработчиков
	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/details", app.details)
	mux.HandleFunc("/model", app.model)
	mux.HandleFunc("/newplant", app.newplant)
	mux.HandleFunc("/about", app.about)
	mux.HandleFunc("/edit", app.editplant)
	mux.HandleFunc("/getplant", app.getplant)


	fileServer := http.FileServer(neuteredFileSystem{http.Dir("./ui/static")})
	mux.Handle("/static", http.NotFoundHandler())
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	return mux
}
