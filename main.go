package main

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var databaseUser string = ""
var databasePassword string = ""
var databaseHost string = ""
var databasePort int32 = 0
var databaseName string = ""

type PGXTestTable struct {
	ColumnOne   int32   `db:"column_one"`
	ColumnTwo   int32   `db:"columntwo"`
	ColumnThree []int32 `db:"column_three`
	ColumnFour  []int32 `db:"columnfour`
}
type Wrapper struct {
	Values []PGXTestTable `db:"values"`
}

func main() {
	if databaseUser == "" {
		fmt.Println("Please add the username for the test database into the main.go file")
		return
	}
	if databasePassword == "" {
		fmt.Println("Please add the password for the test database into the main.go file")
		return
	}
	if databaseHost == "" {
		fmt.Println("Please add the hostname for the test database into the main.go file")
		return
	}
	if databasePort == 0 {
		fmt.Println("Please add the port for the test database into the main.go file")
		return
	}
	if databaseName == "" {
		fmt.Println("Please add the database name for the test database into the main.go file")
		return
	}
	configString := fmt.Sprintf("user=%s password=%s host=%s port=%d dbname=%s", databaseUser, databasePassword, databaseHost, databasePort, databaseName)
	config, err := pgxpool.ParseConfig(configString)
	if err != nil {
		fmt.Println(err)
		return
	}
	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		fmt.Println(err)
		return
	}
	_, err = pool.Exec(context.Background(), `
		CREATE TABLE IF NOT EXISTS pgx_test_table (
			column_one 	 INT,
			columntwo	 INT,
			column_three INT[],
			columnfour   INT[]
		);
	`)
	if err != nil {
		fmt.Println(err)
		return
	}
	_, err = pool.Exec(context.Background(), `
		INSERT INTO pgx_test_table
		(column_one, columntwo, column_three, columnfour)
		VALUES
		(1001, 1011, '{1110}', '{1100}'),
		(2002, 2022, '{2220}', '{2200}');
	`)
	if err != nil {
		fmt.Println(err)
		return
	}
	rows, err := pool.Query(context.Background(), `
		SELECT 
			COALESCE(
				(SELECT 
					json_agg(row_to_json(pgx_test_table))
					FROM pgx_test_table
				),
				'[]'::json
			) values
	`)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer rows.Close()
	wrapper, err := pgx.CollectRows[Wrapper](rows, pgx.RowToStructByName[Wrapper])
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, thing := range wrapper[0].Values {
		fmt.Println("Got values: ", thing.ColumnOne, thing.ColumnTwo, thing.ColumnThree, thing.ColumnFour)
	}
}
