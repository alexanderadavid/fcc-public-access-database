package main

import (
	"database/sql"
	"path/filepath"

	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/go-sql-driver/mysql"
)

func main() {
	db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	initTables(db)
	seedData(db)
}

func seedData(db *sql.DB) {
	path, err := filepath.Abs("config.json")
	if err != nil {
		panic(err)
	}
	config, _ := ioutil.ReadFile(path)
	var files map[string][]string
	err = json.Unmarshal(config, &files)
	if err != nil {
		panic(err)
	}
	for dir, dats := range files {
		for _, dat := range dats {
			uploadDatFile(dir, dat, db)
		}
	}
}

func uploadDatFile(dir, dat string, db *sql.DB) {
	path, err := filepath.Abs(fmt.Sprintf("wirelessftp.fcc.gov/pub/uls/complete/%s/%s.dat", dir, dat))
	if err != nil {
		panic(err)
	}
	fmt.Println("Importing data from " + path)

	_, err = db.Exec("USE WirelessPA;")
	if err != nil {
		panic(err)
	}

	mysql.RegisterLocalFile(path)

	query := fmt.Sprintf(`
	LOAD DATA LOCAL INFILE '%s' IGNORE
	INTO TABLE %s
	COLUMNS TERMINATED BY '|'
	LINES TERMINATED BY '\n';
	`, path, dat)

	_, err = db.Exec(query)
	if err != nil {
		panic(err)
	}
}

func initTables(db *sql.DB) {
	fmt.Println("Initializing Database Tables")
	path, err := filepath.Abs("init.sql")
	if err != nil {
		panic(err)
	}

	c, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	queries := strings.Split(string(c), ";")
	for _, query := range queries[:len(queries)-1] {
		_, err = db.Exec(query + ";")
		if err != nil {
			panic(err)
		}
	}
}
