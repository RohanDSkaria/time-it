package repository

import (
	"database/sql"

	"github.com/RohanDSkaria/time-it/internal/model"
)

type Repository struct {
	DB *sql.DB
}

func New(db *sql.DB) *Repository {
	return &Repository{DB: db}
}

func (r *Repository) GetCurrentEntry() (*model.CurrentEntry, error) {
	row := r.DB.QueryRow("SELECT task, start_time FROM curren_entry WHERE id = 1")

	var c model.CurrentEntry
	err := row.Scan(&c.Task, &c.Start)
	if err != nil {
		return nil, err
	}

	return &c, nil
}

func (r *Repository) SetCurrentEntry(task string, start_time int64) error {
	_, err := r.DB.Exec(
		"INSERT OR REPLACE INTO current_entry(id, task, start_time) VALUES(1,?,?)",
		task, start_time,
	)
	return err
}

func (r *Repository) DeleteCurrentEntry() error {
	_, err := r.DB.Exec("DELETE FROM current_entry WHERE id = 1")
	return err
}

func (r *Repository) AddEntry(e model.Entry) error {
	_, err := r.DB.Exec(
		"INSERT INTO entries(task, start_time, duration) VALUES(?,?,?)",
		e.Task, e.Start, e.Duration,
	)
	return err
}
