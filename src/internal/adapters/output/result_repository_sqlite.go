package output

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"ports-adapters-study/src/internal/core/domain"
)

type ResultStorageSqlite struct {
	db *sql.DB
}

func NewResultStorageSqlite() *ResultStorageSqlite {
	db, err := sql.Open("sqlite3", "store.db")
	if err != nil {
		log.Fatal(err)
	}
	// TODO переместить инит в другое место
	_, err = db.Exec(`CREATE TABLE  IF NOT EXISTS results (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		num INTEGER NOT NULL
	)`)
	if err != nil {
		log.Fatal(err)
	}
	return &ResultStorageSqlite{
		db,
	}
}

func (r *ResultStorageSqlite) AddResult(result domain.ParseResult) error {
	_, err := r.db.Exec("INSERT INTO results (num) VALUES (?)", result.ApsNum)
	if err != nil {
		return err
	}
	return nil
}

func (r *ResultStorageSqlite) GetAllResults() ([]domain.ParseResult, error) {
	rows, err := r.db.Query("SELECT id, num FROM results")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	results := make([]domain.ParseResult, 0)
	for rows.Next() {
		var result domain.ParseResult
		err := rows.Scan(&result.ID, &result.ApsNum)
		if err != nil {
			return nil, err
		}
		results = append(results, result)
	}
	return results, nil
}
