package handler

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"

	"person-service/internal/model"
	"person-service/internal/repository"
	"person-service/internal/service"
)

type Handler struct {
	Repo *repository.PersonRepository
}

func NewHandler(repo *repository.PersonRepository) *Handler {
	return &Handler{Repo: repo}
}

func (h *Handler) CreatePerson(w http.ResponseWriter, r *http.Request) {
	var p model.Person
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	enrichPerson(&p)

	id, err := h.Repo.Create(&p)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Database error")
		return
	}
	p.ID = id

	respondJSON(w, http.StatusOK, p)
}

func (h *Handler) GetPeople(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	limit := parseInt(q.Get("limit"), 10)
	offset := parseInt(q.Get("offset"), 0)

	filters := map[string]interface{}{}
	for _, key := range []string{"name", "surname", "gender"} {
		if val := q.Get(key); val != "" {
			filters[key] = val
		}
	}
	if age := parseInt(q.Get("age"), -1); age >= 0 {
		filters["age"] = age
	}

	people, err := h.Repo.GetAll(filters, limit, offset)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Database error")
		return
	}

	respondJSON(w, http.StatusOK, people)
}

func (h *Handler) UpdatePerson(w http.ResponseWriter, r *http.Request) {
	id := parseInt(mux.Vars(r)["id"], -1)
	if id < 0 {
		respondError(w, http.StatusBadRequest, "Invalid ID")
		return
	}

	var p model.Person
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}
	p.ID = id

	if err := h.Repo.Update(id, &p); err != nil {
		respondError(w, http.StatusInternalServerError, "Database error")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) DeletePerson(w http.ResponseWriter, r *http.Request) {
	id := parseInt(mux.Vars(r)["id"], -1)
	if id < 0 {
		respondError(w, http.StatusBadRequest, "Invalid ID")
		return
	}

	if err := h.Repo.Delete(id); err != nil {
		respondError(w, http.StatusNotFound, "Person not found")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func enrichPerson(p *model.Person) {
	if age, err := service.GetAge(p.Name); err == nil {
		p.Age = age
	}
	if gender, err := service.GetGender(p.Name); err == nil {
		p.Gender = gender
	}
	if nat, err := service.GetNationality(p.Name); err == nil {
		p.Nationality = nat
	}
}

func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}

func respondError(w http.ResponseWriter, status int, msg string) {
	http.Error(w, msg, status)
}

func parseInt(s string, def int) int {
	if v, err := strconv.Atoi(s); err == nil {
		return v
	}
	return def
}
