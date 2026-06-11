package sqlgo

import (
	"database/sql"
	"ecp_app/internal/model"
	"fmt"
)

type SQLRepository struct {
	db *sql.DB
}

func NewSQLRepository(db *sql.DB) (*SQLRepository, error) {
	if db == nil {
		return nil, fmt.Errorf("db is required")
	}

	return &SQLRepository{db: db}, nil
}

func (r *SQLRepository) Add(ecp model.Ecp) error {
	// r.db.Exec(...)

	return nil
}

func (r *SQLRepository) FindExpiring() ([]model.Ecp, error) {
	// r.db.Exec(...)

	return nil, nil
}
