package apiserver

import (
	"encoding/json"
	"errors"
	"home/fosen/Document/golang/RestAPI/internal/app/store"
	"home/fosen/Document/golang/RestAPI/internal/model"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
)

const (
	sessionName = "RestAPI_Auth"
)

var (
	errIncorrectEmailOrPassword = errors.New("incorrect email or password")
)

type server struct {
	router       *mux.Router
	store        store.Store
	sessionStore sessions.Store
}

func newServer(store store.Store, sessionStore sessions.Store) *server {
	s := &server{
		router:       mux.NewRouter(),
		store:        store,
		sessionStore: sessionStore,
	}

	s.configureRouter()
	return s
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *server) configureRouter() {
	s.router.HandleFunc("/users/create", s.handleUsersCreate()).Methods("POST")
	s.router.HandleFunc("/sessions", s.handleSessionsCreate()).Methods("POST")
	s.router.HandleFunc("/task", s.handlerTaskCreate()).Methods("POST")
	s.router.HandleFunc("/users/get", s.handleGetUser()).Methods("POST")
	s.router.HandleFunc("/users", s.handleGetAllUsers()).Methods("POST")
}

func (s *server) handleUsersCreate() http.HandlerFunc {

	type request struct {
		Email    string `json:"email"`
		UserName string `json:"user_name"`
		Password string `json:"password"`
		Role     string `json:"user_role"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		u := &model.User{
			Email:    req.Email,
			UserName: req.UserName,
			Password: req.Password,
			UserRole: req.Role,
		}

		if err := s.store.User().Create(u); err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}

		u.Sanitize()
		s.respond(w, r, http.StatusCreated, u)
	}
}

func (s *server) handleSessionsCreate() http.HandlerFunc {
	type request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		u, err := s.store.User().FindByMail(req.Email)
		if err != nil || !u.ComparePassword(req.Password) {
			s.error(w, r, http.StatusUnauthorized, errIncorrectEmailOrPassword)
			return
		}

		session, err := s.sessionStore.Get(r, sessionName)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		session.Values["user_id"] = u.ID
		if err := s.sessionStore.Save(r, w, session); err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		s.respond(w, r, http.StatusOK, u)
	}
}

func (s *server) handlerTaskCreate() http.HandlerFunc {

	type request struct {
		Email_curator  string `json:"email_curator"`
		Email_employee string `json:"email_employee"`
		Description    string `json:"description"`
		Reward         int    `json:"reward"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		u, err := s.store.User().FindByMail(req.Email_curator)
		if err != nil || u.UserRole != "curator" {
			s.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}

		t := &model.Task{
			Email_curator:  req.Email_curator,
			Email_employee: req.Email_employee,
			Description:    req.Description,
			Reward:         req.Reward,
		}

		if err := s.store.Task().Create(t); err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}

		s.respond(w, r, http.StatusCreated, t)

	}
}

func (s *server) handleGetUser() http.HandlerFunc {

	type request struct {
		Id int `json:"id"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		u, err := s.store.User().FindById(req.Id)
		if err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}

		s.respond(w, r, http.StatusAccepted, u)
	}
}

func (s *server) handleGetAllUsers() http.HandlerFunc {

	type request struct {
		Role string `json:"role"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		if req.Role != "" {
			list_u, err := s.store.User().GetAllUser_filter(req.Role)
			if err != nil {
				s.error(w, r, http.StatusUnprocessableEntity, err)
				return
			}
			s.respond(w, r, http.StatusAccepted, list_u)
			return
		}

		list_u, err := s.store.User().GetAllUser()
		if err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}

		s.respond(w, r, http.StatusAccepted, list_u)
	}
}

func (s *server) handleTaskComplete() http.HandlerFunc {

	type request struct {
		Id    int    `json:"id"`
		Email string `json:"email"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		err1 := s.store.Task().StatusUpdate(req.Email)
		if err1 != nil {
			s.error(w, r, http.StatusUnprocessableEntity, err1)
			return
		}
		s.respond(w, r, http.StatusAccepted, "Updated")

		err2 := s.store.User().LevelUpdate(req.Email, req.Id)
		if err2 != nil {
			s.error(w, r, http.StatusUnprocessableEntity, err2)
			return
		}
		s.respond(w, r, http.StatusAccepted, "Updated")
	}
}

func (s *server) error(w http.ResponseWriter, r *http.Request, code int, err error) {
	s.respond(w, r, code, map[string]string{"error": err.Error()})
}

func (s *server) respond(w http.ResponseWriter, r *http.Request, code int, data interface{}) {
	w.WriteHeader(code)
	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}
