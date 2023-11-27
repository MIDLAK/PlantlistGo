package main

import (
	"fmt"
	"net/http"
	"runtime/debug"
)

//помощник, записывающий сообшение об ошибке в errorLog и отправляющий
//пользователю ошибку 500: "Внутренняя ошибка сервера"
func (app *application) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	//установка глубины вызова на 2 и печать
	app.errorLog.Output(2, trace)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

//помощник,  отправляющий определенный код состояния и соответствующее описание
func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

//помощник, отправляющий "404 Страница не найдена"
func (app *application) notFound(w http.ResponseWriter) {
	app.clientError(w, http.StatusNotFound)
}
