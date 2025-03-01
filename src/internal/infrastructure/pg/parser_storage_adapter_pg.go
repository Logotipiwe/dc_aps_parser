package pg

import (
	"database/sql"
	"dc-aps-parser/src/internal/core/domain"
)

type ParserStorageAdapterPg struct {
	db *sql.DB
}

func NewParserStorageAdapterPg(db *sql.DB) *ParserStorageAdapterPg {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS parsers (
		id BIGINT PRIMARY KEY,
		link VARCHAR NOT NULL,
		username VARCHAR NOT NULL
	);`)
	if err != nil {
		panic(err)
	}
	return &ParserStorageAdapterPg{
		db: db,
	}
}

func (p *ParserStorageAdapterPg) SaveParser(parserData domain.ParserData) error {
	query := `INSERT INTO parsers (id, link, username)
	VALUES ($1, $2, $3)
	ON CONFLICT (id) DO NOTHING;`

	_, err := p.db.Exec(query, parserData.ChatID, parserData.Link, parserData.UserName)
	if err != nil {
		return err
	}
	return nil
}

func (p *ParserStorageAdapterPg) RemoveParser(chatID int64) error {
	query := `DELETE FROM parsers WHERE id = $1;`
	_, err := p.db.Exec(query, chatID)
	if err != nil {
		return err
	}
	return nil
}

func (p *ParserStorageAdapterPg) GetParsers() ([]domain.ParserData, error) {
	rows, err := p.db.Query(`SELECT id, link, username FROM parsers;`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	parsers := make([]domain.ParserData, 0)
	for rows.Next() {
		parserData := domain.ParserData{}
		err := rows.Scan(&parserData.ChatID, &parserData.Link, &parserData.UserName)
		if err != nil {
			return nil, err
		}
		parsers = append(parsers, parserData)
	}
	return parsers, nil
}
