package pg

import (
	"database/sql"
	"encoding/json"
)

type ResultStorageAdapterPg struct {
	db *sql.DB
}

func NewResultStorageAdapterPg(db *sql.DB) *ResultStorageAdapterPg {
	return &ResultStorageAdapterPg{
		db: db,
	}
}

func (r *ResultStorageAdapterPg) SaveNewRawItem(apId int64, rawItem map[string]interface{}) error {
	jsonData, err := json.Marshal(rawItem)
	if err != nil {
		return err
	}
	_, err = r.db.Exec(`INSERT INTO results_history (ap_id, data) VALUES ($1,$2)
		ON CONFLICT (ap_id) DO UPDATE SET data = excluded.data
	`, apId, string(jsonData))
	if err != nil {
		return err
	}
	return nil
}
