package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
    _ "github.com/go-sql-driver/mysql"
    "database/sql"
	"GoTest/pkg/models/mysql"
)

//зависимости приложения
type application struct {
	errorLog 		*log.Logger
	infoLog  		*log.Logger
	dbPlantlist		*mysql.PlantlistModel
}

func main() {
	//флаги командной строки
	addr := flag.String("addr", ":3000", "Сетевой адрес HTTP")
    dbAddr := flag.String("dbAddr", "root:@tcp(127.0.0.1:3306)/plantlist_db", 
                                                            "Сетевой адрес базы данных")
	dbDriverName := flag.String("dbName", "mysql", "Название базы данных (по умолчанию mysql)")
	flag.Parse()

	//создание логгеров
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	//создание соединения с БД и проверка подключения
	db, err := openDB(*dbDriverName, *dbAddr)
	if err != nil {
		errorLog.Fatal(err)
	}

	//срок жизни соединения в пуле с момента его создания
	db.SetConnMaxLifetime(time.Minute * 3)
	//максимальное количество открытых соединений с БД из пула
	db.SetMaxOpenConns(10) 
	/*максимальное количнство соединений, которые могут храниться в пуле и
	ожидать использования*/
	db.SetMaxIdleConns(10)

	//отложенное закрытие соединения (после завершения функции)
	defer db.Close()

    infoLog.Print("Соединение с базой данных установлено")

	//инициализация структуры с зависимостями приложения
	app := &application{
		errorLog:	 errorLog,
		infoLog:  	 infoLog,
		dbPlantlist: &mysql.PlantlistModel{DB: db},
	}

	//инициализация новвй структуры http.Server
	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  app.routes(), //получение маршрутизатора
	}

	infoLog.Printf("Запуск веб-сервера на %v", *addr)
	srvErr := srv.ListenAndServe()
	errorLog.Fatal(srvErr)

}

type neuteredFileSystem struct {
	fs http.FileSystem
}

//вызвается каждый раз, когда FileServer получает запрос
//и защищает от прямого перехода в папку со статическими файлами
func (nfs neuteredFileSystem) Open(path string) (http.File, error) {
	f, err := nfs.fs.Open(path)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	s, _ := f.Stat()
	if s.IsDir() {
		index := filepath.Join(path, "index.html")
		if _, err := nfs.fs.Open(index); err != nil {
			closeErr := f.Close()
			if closeErr != nil {
				return nil, closeErr
			}

			return nil, err
		}
	}

	return f, nil
}

//возвращает пул соединений sql.DB для заданной строки подключения DNS
func openDB(driver string, dns string) (*sql.DB, error) {
	//инициализация пула подключений
	db, err := sql.Open(driver, dns) 
	if err != nil {
		return nil, err
	}

	//создание соединения для проверки
	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
