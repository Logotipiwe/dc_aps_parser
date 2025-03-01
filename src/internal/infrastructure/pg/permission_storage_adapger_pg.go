package pg

import (
	"database/sql"
	"errors"
	"log"
)

type PermissionStorageAdapterPg struct {
	db *sql.DB
}

func NewPermissionStorageAdapterPg(db *sql.DB) *PermissionStorageAdapterPg {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS permissions (
    	chat_id bigint NOT NULL PRIMARY KEY,
    	permitted_aps_num integer NOT NULL
	)`)
	if err != nil {
		log.Fatal(err)
	}
	return &PermissionStorageAdapterPg{db: db}
}

func (p *PermissionStorageAdapterPg) GetPermittedApsNumForChat(chatID int64) (*int, error) {
	row := p.db.QueryRow("SELECT permitted_aps_num FROM permissions WHERE chat_id = $1", chatID)
	var permittedApsNum int
	err := row.Scan(&permittedApsNum)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &permittedApsNum, nil
}
