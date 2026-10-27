package main

import (
	"database/sql"
	"time"

	_ "modernc.org/sqlite"
)

func openDB(path string) (*sql.DB, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}
	schema := `
CREATE TABLE IF NOT EXISTS enquiries (
  id          INTEGER PRIMARY KEY AUTOINCREMENT,
  name        TEXT NOT NULL,
  email       TEXT NOT NULL,
  type        TEXT NOT NULL,
  message     TEXT NOT NULL,
  status      TEXT NOT NULL DEFAULT 'New',
  created_at  DATETIME DEFAULT CURRENT_TIMESTAMP
);`
	if _, err := db.Exec(schema); err != nil {
		return nil, err
	}
	return db, nil
}

type Enquiry struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Email       string    `json:"email"`
	Type        string    `json:"type"`
	Message     string    `json:"message"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	FollowUpDue bool      `json:"follow_up_due"`
}

func insertEnquiry(db *sql.DB, name, email, enquiryType, message string) (int64, error) {
	res, err := db.Exec(`INSERT INTO enquiries (name, email, type, message) VALUES (?, ?, ?, ?)`,
		name, email, enquiryType, message)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func listEnquiry(db *sql.DB, id int64) (*Enquiry, error) {
	row := db.QueryRow(`
		SELECT id, name, email, type, message, status, created_at
		FROM enquiries
		WHERE id = ?`, id)
	var e Enquiry
	var created string
	if err := row.Scan(&e.ID, &e.Name, &e.Email, &e.Type, &e.Message, &e.Status, &created); err != nil {
		return nil, err
	}
	e.CreatedAt = parseSQLiteTime(created)
	return &e, nil
}

func listEnquiries(db *sql.DB) ([]Enquiry, error) {
	rows, err := db.Query(`
		SELECT id, name, email, type, message, status, created_at
		FROM enquiries
		ORDER BY created_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	cutoff := time.Now().UTC().Add(-48 * time.Hour)
	list := []Enquiry{}
	for rows.Next() {
		var e Enquiry
		var created string
		if err := rows.Scan(&e.ID, &e.Name, &e.Email, &e.Type, &e.Message, &e.Status, &created); err != nil {
			return nil, err
		}
		e.CreatedAt = parseSQLiteTime(created)
		e.FollowUpDue = e.Status == "New" && e.CreatedAt.Before(cutoff)
		list = append(list, e)
	}
	return list, rows.Err()
}

func updateStatus(db *sql.DB, id int64, status string) (bool, error) {
	res, err := db.Exec(`UPDATE enquiries SET status = ? WHERE id = ?`, status, id)
	if err != nil {
		return false, err
	}
	n, err := res.RowsAffected()
	return n != 0, err
}

func backdateEnquiry(db *sql.DB, id int64) (bool, error) {
	res, err := db.Exec(
		`UPDATE enquiries SET created_at = datetime('now', '-3 days') WHERE id = ?`,
		id,
	)
	if err != nil {
		return false, err
	}
	n, err := res.RowsAffected()
	return n != 0, err
}

func parseSQLiteTime(s string) time.Time {
	for _, layout := range []string{time.RFC3339, "2006-01-02 15:04:05"} {
		if t, err := time.ParseInLocation(layout, s, time.UTC); err == nil {
			return t
		}
	}
	return time.Now().UTC()
}
