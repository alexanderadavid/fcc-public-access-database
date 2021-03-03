package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"testing"
)

func Test_buildInsertQuery(t *testing.T) {
	type args struct {
		dat  string
		rows []string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			"builds an insert query",
			args{
				"table",
				[]string{"1|2|3|4|5", "6|7|8|9|10"},
			},
			"insert ignore into table values ('1', '2', '3', '4', '5'), ('6', '7', '8', '9', '10')",
		},
		{
			"builds an insert query with blank values",
			args{
				"table",
				[]string{"||||", "||||"},
			},
			"insert ignore into table values ('', '', '', '', ''), ('', '', '', '', '')",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ch := make(chan string, len(tt.args.rows))
			go func() {
				defer close(ch)
				for _, row := range tt.args.rows {
					ch <- row
				}
			}()

			if got := buildInsertQuery(tt.args.dat, ch); got != tt.want {
				t.Errorf("buildInsertQuery() = %v, want %v", got, tt.want)
			}
		})
	}
	return
}

func Test_initTables(t *testing.T) {
	want := []string{
		"create database if not exists WirelessPA;",
		"use WirelessPA;",
	}

	queryCount := 0
	t.Run("initializes tables", func(t *testing.T) {
		exec := func(query string, args ...interface{}) (sql.Result, error) {
			// only checking for the shorter queries in init.sql to reduce LOC
			if queryCount >= len(want) {
				return nil, nil
			}
			if want[queryCount] != query {
				return nil, fmt.Errorf("%s should equal %s", want[queryCount], query)
			}
			queryCount++
			return nil, nil
		}
		if err := initTables(exec); err != nil {
			t.Errorf("initTables() error = %v", err)
		}
	})
}

func Test_uploadDatFile(t *testing.T) {
	called := []string{}
	type args struct {
		table string
		exec  execType
	}
	tests := []struct {
		name         string
		args         args
		fileContents []byte
		rowBatchSize int
		numWorkers   int
		want         []string
		wantErr      bool
	}{
		{
			"uploads a dat file with one worker (deterministic)",
			args{
				"testTable",
				func(query string, args ...interface{}) (sql.Result, error) {
					called = append(called, query)
					return nil, nil
				},
			},
			[]byte(`1|1|1|1
2|2|2|2
3|3|3|3
4|4|4|4
5|5|5|5
6|6|6|6
7|7|7|7
8|8|8|8
9|9|9|9`),
			4,
			1,
			[]string{
				"insert ignore into testTable values ('1', '1', '1', '1'), ('2', '2', '2', '2'), ('3', '3', '3', '3'), ('4', '4', '4', '4')",
				"insert ignore into testTable values ('5', '5', '5', '5'), ('6', '6', '6', '6'), ('7', '7', '7', '7'), ('8', '8', '8', '8')",
				"insert ignore into testTable values ('9', '9', '9', '9')",
			},
			false,
		},
	}
	for _, tt := range tests {
		// reset configuration vars to  reasonable testing values
		rowBatchSize = tt.rowBatchSize
		numWorkers = tt.numWorkers
		// create in memory file system
		dir, err := ioutil.TempDir("", "")
		if err != nil {
			log.Fatal(err)
		}
		defer os.RemoveAll(dir)

		tmpfn := filepath.Join(dir, "tmpfile.dat")
		if err := ioutil.WriteFile(tmpfn, tt.fileContents, 0666); err != nil {
			log.Fatal(err)
		}

		t.Run(tt.name, func(t *testing.T) {
			if err := uploadDatFile(tmpfn, tt.args.table, tt.args.exec); (err != nil) != tt.wantErr {
				t.Errorf("uploadDatFile() error = %v, wantErr %v", err, tt.wantErr)
			}
			for i, call := range tt.want {
				if call != called[i] {
					t.Errorf("uploadDatFile() query %v, want %v", called[i], call)
				}
			}
			called = []string{}
		})
	}
}
