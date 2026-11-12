package db

import (
	"database/sql"
	"ecp_app/internal/model"
	"os"
	"path/filepath"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type Database struct {
	DB *sql.DB
}

func Init() (*Database, error) {
	exePath, err := os.Executable()
	if err != nil {
		return nil, err
	}
	dbPath := filepath.Join(filepath.Dir(exePath), "app.Databa")
	dsn := dbPath + "?_busy_timeout=5000&_foreign_keys=on"

	conn, err := sql.Open("sqlite3", dsn)
	if err != nil {
		return nil, err
	}
	if err := conn.Ping(); err != nil {
		conn.Close()
		return nil, err
	}

	_, err = conn.Exec(`
		CREATE TABLE IF NOT EXISTS certificates (
			id          INTEGER PRIMARY KEY AUTOINCREMENT,
			second_name TEXT NOT NULL,
			first_name  TEXT NOT NULL,
			patronymic  TEXT NOT NULL,
			district    TEXT,
			expiration  DATETIME NOT NULL,
			UNIQUE(second_name, first_name, patronymic)
		);

		CREATE TABLE IF NOT EXISTS districts (
			id   INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL UNIQUE
		);
	`)
	if err != nil {
		conn.Close()
		return nil, err
	}

	if err := seedDistrictsIfEmpty(conn); err != nil {
		conn.Close()
		return nil, err
	}

	return &Database{DB: conn}, nil
}

func (db *Database) SaveCert(cer model.Certificate) error {
	_, err := db.DB.Exec(`
		INSERT INTO certificates (second_name, first_name, patronymic, district, expiration)
		VALUES (?, ?, ?, ?, ?)
		ON CONFLICT(second_name, first_name, patronymic) DO UPDATE SET
			expiration = excluded.expiration;
	`, cer.SecondName, cer.FirstName, cer.Patronymic, cer.District, cer.Expiration)
	return err
}

func (db *Database) LoadCerts() ([]model.Certificate, error) {
	rows, err := db.DB.Query(`
		SELECT second_name, first_name, patronymic, district, expiration
		FROM certificates;
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make([]model.Certificate, 0)
	for rows.Next() {
		var c model.Certificate
		if err := rows.Scan(
			&c.SecondName,
			&c.FirstName,
			&c.Patronymic,
			&c.District,
			&c.Expiration,
		); err != nil {
			return nil, err
		}
		result = append(result, c)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

func (db *Database) LoadDistricts() ([]string, error) {
	rows, err := db.DB.Query(`SELECT name FROM districts ORDER BY name;`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, err
		}
		result = append(result, name)
	}
	return result, rows.Err()
}

func (db *Database) LoadExpiringNextMonth() ([]model.Certificate, error) {
	// Вычисляем границы следующего месяца
	now := time.Now()
	firstOfNext := time.Date(now.Year(), now.Month()+1, 1, 0, 0, 0, 0, now.Location())
	firstOfAfter := firstOfNext.AddDate(0, 1, 0)

	rows, err := db.DB.Query(`
		SELECT second_name, first_name, patronymic, district, expiration
		FROM certificates
		WHERE expiration >= ? AND expiration < ?
		ORDER BY expiration ASC;
	`, firstOfNext, firstOfAfter)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make([]model.Certificate, 0)
	for rows.Next() {
		var c model.Certificate
		if err := rows.Scan(
			&c.SecondName,
			&c.FirstName,
			&c.Patronymic,
			&c.District,
			&c.Expiration,
		); err != nil {
			return nil, err
		}
		result = append(result, c)
	}
	return result, rows.Err()
}

func seedDistrictsIfEmpty(conn *sql.DB) error {
	if len(seedDistricts) == 0 {
		return nil
	}

	tx, err := conn.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare(`INSERT OR IGNORE INTO districts (name) VALUES (?);`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, name := range seedDistricts {
		if _, err := stmt.Exec(name); err != nil {
			return err
		}
	}

	return tx.Commit()
}
