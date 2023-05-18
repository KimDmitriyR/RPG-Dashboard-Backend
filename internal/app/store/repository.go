package store

import "home/fosen/Document/golang/RestAPI/internal/model"

type UserRepository interface {
	Create(*model.User) error
	FindByMail(string) (*model.User, error)
	FindById(int) (*model.User, error)
}

type TaskRepository interface {
	Create(*model.Task) error
}
