package store

import (
	"database/sql"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

// Repo represents a GitHub repository.
type Repo struct {
	Name      string
	Stars     int
	Language  string
	UpdatedAt string
	URL       string
}

// Job represents a queued job.
type Job struct {
	ID      int64
	Cmd     string
	Status  string
	Retries int
}

// Store wraps the sqlite DB.
type Store struct {
	db *sql.DB
}

// Open opens or creates the sqlite DB at path, ensuring WAL mode and required tables.
func Open(path string) (*Store, error) {
	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		return nil, err
	}
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}
	if _, err := db.Exec("PRAGMA journal_mode=WAL"); err != nil {
		db.Close()
		return nil, err
	}
	if _, err := db.Exec(`CREATE TABLE IF NOT EXISTS repos(
		name TEXT,
		stars INTEGER,
		language TEXT,
		updated_at TEXT,
		url TEXT PRIMARY KEY
	)`); err != nil {
		db.Close()
		return nil, err
	}
	if _, err := db.Exec(`CREATE TABLE IF NOT EXISTS jobs(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		cmd TEXT,
		status TEXT,
		retries INTEGER
	)`); err != nil {
		db.Close()
		return nil, err
	}
	return &Store{db: db}, nil
}

// UpsertRepos inserts or replaces repos. If URL is empty, Name is used as URL to avoid PK collision.
func (s *Store) UpsertRepos(repos []Repo) error {
	for _, r := range repos {
		url := r.URL
		if url == "" {
			url = r.Name
		}
		_, err := s.db.Exec(
			`INSERT OR REPLACE INTO repos(name, stars, language, updated_at, url) VALUES(?,?,?,?,?)`,
			r.Name, r.Stars, r.Language, r.UpdatedAt, url,
		)
		if err != nil {
			return err
		}
	}
	return nil
}

// ListRepos returns repos sorted by sortBy ("stars" -> stars DESC, otherwise name ASC).
func (s *Store) ListRepos(sortBy string) ([]Repo, error) {
	order := "stars DESC"
	switch sortBy {
	case "stars":
		order = "stars DESC"
	case "name":
		order = "name ASC"
	case "updated":
		order = "updated_at DESC"
	default:
		// default to stars for now
		order = "stars DESC"
	}
	rows, err := s.db.Query(`SELECT name, stars, language, updated_at, url FROM repos ORDER BY ` + order)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []Repo
	for rows.Next() {
		var r Repo
		if err := rows.Scan(&r.Name, &r.Stars, &r.Language, &r.UpdatedAt, &r.URL); err != nil {
			return nil, err
		}
		out = append(out, r)
	}
	return out, rows.Err()
}

// PushJob inserts a pending job and returns its id.
func (s *Store) PushJob(cmd string) (int64, error) {
	res, err := s.db.Exec(`INSERT INTO jobs(cmd, status, retries) VALUES(?,?,?)`, cmd, "pending", 0)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

// ListJobs returns jobs filtered by status. If status is empty, returns all.
func (s *Store) ListJobs(status string) ([]Job, error) {
	var rows *sql.Rows
	var err error
	if status == "" {
		rows, err = s.db.Query(`SELECT id, cmd, status, retries FROM jobs ORDER BY id ASC`)
	} else {
		rows, err = s.db.Query(`SELECT id, cmd, status, retries FROM jobs WHERE status=? ORDER BY id ASC`, status)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []Job
	for rows.Next() {
		var j Job
		if err := rows.Scan(&j.ID, &j.Cmd, &j.Status, &j.Retries); err != nil {
			return nil, err
		}
		out = append(out, j)
	}
	return out, rows.Err()
}

// Close closes the database.
func (s *Store) Close() error {
	return s.db.Close()
}

func (s *Store) SetJobStatus(id int64, status string) error {
	_, err := s.db.Exec(`UPDATE jobs SET status=? WHERE id=?`, status, id)
	return err
}

func (s *Store) IncrementRetries(id int64) error {
	_, err := s.db.Exec(`UPDATE jobs SET retries=retries+1 WHERE id=?`, id)
	return err
}

func (s *Store) ListTables() ([]string, error) {
	rows, err := s.db.Query("SELECT name FROM sqlite_master WHERE type='table' AND name NOT LIKE 'sqlite_%' ORDER BY name")
	if err != nil { return nil, err }
	defer rows.Close()
	var out []string
	for rows.Next() {
		var n string
		if err := rows.Scan(&n); err != nil { return nil, err }
		out = append(out, n)
	}
	return out, rows.Err()
}
