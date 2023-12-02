package main

import (
	"html/template"
	"net/http"
	"path/filepath"
	"fmt"
)

//отображение шаблонов
func render[T any](w http.ResponseWriter, r *http.Request, name string, td *T, cache map[string]*template.Template) error {
	//извлечение страниц по названию из кэша
	ts, ok := cache[name]
	if !ok {
		return fmt.Errorf("Шаблон %s не существует или недоступен!", name)
	}

	//рендер файлов шаблона
	err := ts.Execute(w, td)
	if err != nil {
		return err
	}

	return nil
}

//создание кэша шаблонов
func newTemplateCache(dir string) (map[string]*template.Template, error) {
	//хранилище кэша
	cache := map[string]*template.Template{}

	//получение всех файлов шаблонов приложения
	pages, err := filepath.Glob(filepath.Join(dir, "*.page.html"))
	if err != nil {
		return nil, err
	}

	//перебор всех фалов шаблонов
	for _, page := range pages {
		//полный путь файла
		name := filepath.Base(page)

		//файлы шаблонов
		ts, err := template.ParseFiles(page)
		if err != nil {
			return nil, err
		}

		//каркасные шаблоны
		ts, err = ts.ParseGlob(filepath.Join(dir, "*.layout.html"))
		if err != nil {
			return nil, err
		}

		//вспомогательные шаблоны
		ts, err = ts.ParseGlob(filepath.Join(dir, "*.partial.html"))
		if err != nil {
			return nil, err
		} 

		//запись набора шаблонов в кэш (доступ по названию страницы)
		cache[name] = ts
	}

	return cache, nil
}