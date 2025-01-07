package output

import (
	"database/sql"
	"dc-aps-parser/src/internal/core/domain"
)

type ParserStorageAdapterSqlite struct {
	db *sql.DB
}

func NewParserStorageAdapterSqlite(db *sql.DB) *ParserStorageAdapterSqlite {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS parsers (
		id BIGINT PRIMARY KEY,
		link VARCHAR NOT NULL
	);`)
	if err != nil {
		panic(err)
	}
	return &ParserStorageAdapterSqlite{
		db: db,
	}
}

func (p *ParserStorageAdapterSqlite) SaveParser(parserData domain.ParserData) error {
	query := `INSERT INTO parsers (id, link)
	VALUES ($1, $2)
	ON CONFLICT (id) DO NOTHING;`

	_, err := p.db.Exec(query, parserData.ChatID, parserData.Link)
	if err != nil {
		return err
	}
	return nil
}

func (p *ParserStorageAdapterSqlite) RemoveParser(parserData domain.ParserData) error {
	query := `DELETE FROM parsers WHERE id = $1;`
	_, err := p.db.Exec(query, parserData.ChatID)
	if err != nil {
		return err
	}
	return nil
}

func (p *ParserStorageAdapterSqlite) GetParsers() ([]domain.ParserData, error) {
	rows, err := p.db.Query(`SELECT id, link FROM parsers;`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	parsers := make([]domain.ParserData, 0)
	for rows.Next() {
		parserData := domain.ParserData{}
		err := rows.Scan(&parserData.ChatID, &parserData.Link)
		if err != nil {
			return nil, err
		}
		parsers = append(parsers, parserData)
	}
	return parsers, nil
}
