package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"database/sql"
	"encoding/json"
	"io/ioutil"
	"path/filepath"
	"sync/atomic"

	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/sync/errgroup"
)

var rowBatchSize = 50
var numWorkers = 10

const wirelessFtpDir = "wirelessftp.fcc.gov/pub/uls/complete/%s/%s.dat"

// stub db.Exec for unit testing which queries are called
type execType = func(query string, args ...interface{}) (sql.Result, error)

func main() {
	start := time.Now()
	// Open a connection with no database so we can create one
	db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/")
	if err != nil {
		panic(err)
	}
	if err := initTables(db.Exec); err != nil {
		panic(err)
	}
	// Open a connection with the newly created Database
	db, err = sql.Open("mysql", "root:root@/WirelessPA")
	if err != nil {
		panic(err)
	}
	defer db.Close()
	if err := seedData(db.Exec); err != nil {
		panic(err)
	}
	fmt.Println("Finished... took " + time.Since(start).String())
}

func initTables(exec execType) error {
	fmt.Println("Initializing Database Tables")
	path, err := filepath.Abs("init.sql")
	if err != nil {
		return err
	}

	c, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	queries := strings.Split(string(c), ";")
	for _, query := range queries[:len(queries)-1] {
		_, err = exec(strings.TrimSpace(query) + ";")
		if err != nil {
			return err
		}
	}
	return nil
}

func seedData(exec execType) error {
	path, err := filepath.Abs("config.json")
	if err != nil {
		return err
	}
	config, _ := ioutil.ReadFile(path)
	var files map[string][]string
	err = json.Unmarshal(config, &files)
	if err != nil {
		return err
	}
	for dir, dats := range files {
		for _, dat := range dats {
			if err := uploadDatFile(
				fmt.Sprintf(wirelessFtpDir, dir, dat), dat, exec,
			); err != nil {
				return err
			}
		}
	}
	return nil
}

func uploadDatFile(filePath, table string, exec execType) error {
	path, err := filepath.Abs(filePath)
	if err != nil {
		return err
	}
	fmt.Println("Importing data from " + path)

	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	rows := make(chan string)

	g, _ := errgroup.WithContext(context.TODO())

	// atomic bool to sync when we're done
	var finished int32

	// fanning out consumers to run insert queries
	for i := 0; i < numWorkers; i++ {
		g.Go(func() error {
			if err := queryBuilder(table, rows, &finished, exec); err != nil {
				return err
			}
			return nil
		})
	}

	// open the pipeline
	g.Go(func() error {
		defer close(rows)
		for scanner.Scan() {
			rows <- scanner.Text()
		}

		return scanner.Err()
	})
	return g.Wait()
}

func queryBuilder(table string, rows <-chan string, finished *int32, exec execType) error {
	for atomic.LoadInt32(finished) == 0 {
		q := buildInsertQuery(table, rows)
		if q == "" {
			atomic.StoreInt32(finished, 1)
			return nil
		}
		if _, err := exec(q); err != nil {
			fmt.Println("Error executing batch insert query:", err)
			return err
		}
	}
	return nil
}

func buildInsertQuery(table string, rows <-chan string) string {
	insertLine := fmt.Sprintf("insert ignore into %s values ", table)
	values := ""
	rowCount := 0
	for row := range rows {
		row = strings.ReplaceAll(row, "'", "\\'")
		values += fmt.Sprintf("('%s'), ", strings.ReplaceAll(row, "|", "', '"))
		if rowCount >= rowBatchSize - 1 {
			values = values[:len(values)-2] // remove trailing comma and space
			return insertLine + values
		}
		rowCount++
	}
	if rowCount == 0 {
		return ""
	}
	values = values[:len(values)-2] // remove trailing comma and space
	return insertLine + values
}
