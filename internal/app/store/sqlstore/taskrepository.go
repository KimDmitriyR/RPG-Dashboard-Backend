package sqlstore

import "home/fosen/Document/golang/RestAPI/internal/model"

type TaskRepository struct {
	store *Store
}

func (r *TaskRepository) Create(t *model.Task) error {
	return r.store.db.QueryRow(
		"INSERT INTO tasks(email_curator, email_employee, description) values($1, $2, $3) returning id",
		t.Email_curator,
		t.Email_employee,
		t.Description,
	).Scan(&t.ID)
}
