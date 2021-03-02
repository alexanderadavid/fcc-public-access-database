package main

import (
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
