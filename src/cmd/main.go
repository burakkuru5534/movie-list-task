package main

import (
	"encoding/csv"
	"errors"
	"example.com/m/v2/src/api"
	"example.com/m/v2/src/helper"
	"example.com/m/v2/src/service"
	_ "github.com/lib/pq"
	"log"
	"os"
	"strings"
)

func main() {

	conInfo := helper.PgConnectionInfo{
		Host:     "127.0.0.1",
		Port:     5432,
		Database: "movietask",
		Username: "postgres",
		Password: "tayitkan",
		SSLMode:  "disable",
	}

	db, err := helper.NewPgSqlxDbHandle(conInfo, 10)
	if err != nil {
		errors.New("create db handle error.")
	}
	err = db.Ping()
	if err != nil {
		errors.New("ping db error.")
	}

	// Create Appplication Service
	err = helper.InitApp(db)
	if err != nil {
		errors.New("init app error.")
	}
	// open file
	f, err := os.Open("/home/burak/movie-list-task/src/cmd/movies.csv")
	if err != nil {
		log.Fatal(err)
	}

	// remember to close the file at the end of the program
	defer f.Close()

	// read csv values using csv.Reader
	csvReader := csv.NewReader(f)
	data, err := csvReader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	initTable(db, data)

	service.StartHttpService(8000, api.HttpService())
}

func initTable(db *helper.DbHandle, data [][]string) {
	// Create Table
	_, err := db.Exec("CREATE TABLE IF NOT EXISTS movie (id SERIAL PRIMARY KEY, title TEXT unique, time bigint)")
	if err != nil {
		log.Fatal(err)
	}

	// Insert Data
	for i, _ := range data {
		if i == 0 {
			continue
		}
		splittedData := strings.Split(data[i][0], "	")
		title := splittedData[0]
		time := helper.StrToInt64(splittedData[1])
		_, err := db.Exec("INSERT INTO movie (title, time) VALUES ($1, $2) on conflict (title) do nothing", title, time)
		if err != nil {
			log.Fatal(err)
		}
	}
}
