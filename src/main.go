package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	_ "github.com/lib/pq"
)

type DBInterface interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
	QueryRow(query string, args ...interface{}) *sql.Row
	Close() error
}

var db DBInterface

func main() {
	db = initDB()
	defer db.Close()

	http.HandleFunc("/push", pushHandler)
	http.HandleFunc("/avg", avgHandler)

	log.Fatal(http.ListenAndServe(":80", nil))
}

func initDB() *sql.DB {
	connStr := fmt.Sprintf("postgresql://obvallata:sber@app_db/metrics?sslmode=disable")

	database, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("sql.Open error: %v\n", err)
	}

	ensureTableExists(database)
	return database
}

func ensureTableExists(db *sql.DB) {
	_, err := db.Exec(`
        CREATE TABLE IF NOT EXISTS t_metrics (
            time TIMESTAMP,
            metric INTEGER
        );
    `)

	if err != nil {
		log.Fatalf("db.Exec create table error: %v\n", err)
	}
}

func pushHandler(w http.ResponseWriter, r *http.Request) {
	valueStr := r.URL.Query().Get("value")
	if valueStr == "" {
		http.Error(w, `missing value query parameter`, http.StatusBadRequest)
		return
	}

	value, err := strconv.Atoi(valueStr)
	if err != nil {
		http.Error(w, `value should be integer`, http.StatusBadRequest)
		return
	}

	_, err = db.Exec("INSERT INTO t_metrics (time, metric) VALUES ($1, $2)", time.Now(), value)
	if err != nil {
		http.Error(w, fmt.Sprintf("db.Exec INSERT INTO: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Success\n"))
}

func avgHandler(w http.ResponseWriter, r *http.Request) {
	var avg sql.NullFloat64
	err := db.QueryRow("SELECT AVG(metric) FROM t_metrics").Scan(&avg)
	if err != nil {
		http.Error(w, fmt.Sprintf("db.QueryRow SELECT AVG: %v", err), http.StatusInternalServerError)
		return
	}

	if avg.Valid {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprintf("Average: %f\n", avg.Float64)))
	} else {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("No data yet\n"))
	}
}
