package sqlstore

import (
	"database/sql"
	"home/fosen/Document/golang/RestAPI/internal/app/store"
	"home/fosen/Document/golang/RestAPI/internal/model"
)

type TaskRepository struct {
	store *Store
}

func (r *TaskRepository) Create(t *model.Task) error {
	t.Status = "false"
	return r.store.db.QueryRow(
		"INSERT INTO tasks(email_curator, email_employee, description) values($1, $2, $3) returning id",
		t.Email_curator,
		t.Email_employee,
		t.Description,
	).Scan(&t.ID)
}

func (r *TaskRepository) StatusUpdate(email string) error {
	if _, err := r.store.db.Exec("UPDATE tasks set status = $1 FROM tasks WHERE email_empoyee = $2",
		true, email,
	); err != nil {
		if err == sql.ErrNoRows {
			return store.ErrRecordNotFound
		}
		return err
	}
	return nil
}
