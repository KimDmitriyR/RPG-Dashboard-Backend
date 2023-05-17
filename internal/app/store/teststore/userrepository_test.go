package teststore_test

import (
	"home/fosen/Document/golang/RestAPI/internal/app/store"
	"home/fosen/Document/golang/RestAPI/internal/app/store/teststore"
	"home/fosen/Document/golang/RestAPI/internal/model"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUserRepository_Create(t *testing.T) {

	s := teststore.New()
	u := model.TestUser(t)
	assert.NoError(t, s.User().Create(u))
	assert.NotNil(t, u)
}

func TestUserRepository_FindByEmail(t *testing.T) {
	s := teststore.New()
	email := "user@example.org"
	_, err := s.User().FindByMail(email)
	assert.EqualError(t, err, store.ErrRecordNotFound.Error())

	u := model.TestUser(t)
	u.Email = email
	s.User().Create(u)
	u, err = s.User().FindByMail(email)
	assert.NoError(t, err)
	assert.NotNil(t, u)
}
