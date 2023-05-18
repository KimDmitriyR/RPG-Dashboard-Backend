package sqlstore

import (
	"database/sql"
	"home/fosen/Document/golang/RestAPI/internal/app/store"
	"home/fosen/Document/golang/RestAPI/internal/model"
)

type UserRepository struct {
	store *Store
}

func (r *UserRepository) Create(u *model.User) error {

	if err := u.Validate(); err != nil {
		return err
	}

	if err := u.BeforeCreate(); err != nil {
		return err
	}

	return r.store.db.QueryRow(
		"INSERT INTO users(email, name_user, encrypted_password, role) VALUES($1, $2, $3, $4) RETURNING id",
		u.Email,
		u.UserName,
		u.EncryptedPassword,
		u.UserRole,
	).Scan(&u.ID)
}

func (r *UserRepository) FindByMail(email string) (*model.User, error) {
	u := &model.User{}
	if err := r.store.db.QueryRow(
		"SELECT id, email, name_user, encrypted_password, role, user_level FROM users WHERE email = $1",
		email,
	).Scan(
		&u.ID,
		&u.Email,
		&u.UserName,
		&u.EncryptedPassword,
		&u.UserRole,
		&u.UserLevel,
	); err != nil {

		if err == sql.ErrNoRows {
			return nil, store.ErrRecordNotFound
		}

		return nil, err
	}

	return u, nil
}
