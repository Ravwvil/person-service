package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"person-service/internal/logger"
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

// CreatePerson godoc
// @Summary     Создать нового человека
// @Description Создает новую запись с обогащением данных
// @Tags        people
// @Accept      json
// @Produce     json
// @Param       person body     model.CreatePersonRequest true "Данные для создания"
// @Success     200    {object} model.Person
// @Failure     400    {object} model.ErrorResponse
// @Failure     500    {object} model.ErrorResponse
// @Router      /people [post]
func (h *Handler) CreatePerson(w http.ResponseWriter, r *http.Request) {
	var req model.CreatePersonRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Log.Errorw("Invalid JSON", "error", err)
		respondError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}
	logger.Log.Infow("CreatePerson", "request", req)

	p := req.ToPerson()
	if err := enrichPerson(&p); err != nil {
		logger.Log.Warnw("Enrichment warning", "error", err)
	}

	id, err := h.Repo.Create(&p)
	if err != nil {
		logger.Log.Errorw("DB insert failed", "error", err)
		respondError(w, http.StatusInternalServerError, "Database error")
		return
	}
	p.ID = id
	logger.Log.Infow("Person created", "id", id)

	respondJSON(w, http.StatusOK, p)
}

// GetPeople godoc
// @Summary     Получить список людей
// @Description Возвращает всех людей с фильтрами и пагинацией
// @Tags        people
// @Accept      json
// @Produce     json
// @Param       name    query    string false "Имя"
// @Param       surname query    string false "Фамилия"
// @Param       gender  query    string false "Пол"
// @Param       age     query    int    false "Возраст"
// @Param       limit   query    int    false "Лимит"   default(10)
// @Param       offset  query    int    false "Смещение" default(0)
// @Success     200     {array}  model.Person
// @Failure     500     {object} model.ErrorResponse
// @Router      /people [get]
func (h *Handler) GetPeople(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	limit, offset := parseInt(q.Get("limit"), 10), parseInt(q.Get("offset"), 0)

	filters := make(map[string]interface{})
	for _, key := range []string{"name", "surname", "gender"} {
		if v := q.Get(key); v != "" {
			filters[key] = v
		}
	}
	if age := parseInt(q.Get("age"), -1); age >= 0 {
		filters["age"] = age
	}
	logger.Log.Debugw("GetPeople", "filters", filters, "limit", limit, "offset", offset)

	people, err := h.Repo.GetAll(filters, limit, offset)
	if err != nil {
		logger.Log.Errorw("DB select failed", "error", err)
		respondError(w, http.StatusInternalServerError, "Database error")
		return
	}
	logger.Log.Infow("Fetched people", "count", len(people))

	respondJSON(w, http.StatusOK, people)
}

// UpdatePerson godoc
// @Summary     Обновить человека
// @Description Обновляет запись по ID
// @Tags        people
// @Accept      json
// @Produce     json
// @Param       id     path     int                       true "ID"
// @Param       person body     model.UpdatePersonRequest true "Данные для обновления"
// @Success     204    {string} string                   "No Content"
// @Failure     400    {object} model.ErrorResponse
// @Failure     500    {object} model.ErrorResponse
// @Router      /people/{id} [put]
func (h *Handler) UpdatePerson(w http.ResponseWriter, r *http.Request) {
	id := parseInt(mux.Vars(r)["id"], -1)
	if id < 0 {
		respondError(w, http.StatusBadRequest, "Invalid ID")
		return
	}

	var req model.UpdatePersonRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}
	logger.Log.Infow("UpdatePerson", "id", id, "body", req)

	p := req.ToPerson(id)
	if err := h.Repo.Update(id, &p); err != nil {
		logger.Log.Errorw("DB update failed", "error", err)
		respondError(w, http.StatusInternalServerError, "Database error")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// DeletePerson godoc
// @Summary     Удалить человека
// @Description Удаляет запись по ID
// @Tags        people
// @Accept      json
// @Produce     json
// @Param       id   path     int true "ID"
// @Success     204  {string} string           "No Content"
// @Failure     400  {object} model.ErrorResponse
// @Failure     404  {object} model.ErrorResponse
// @Router      /people/{id} [delete]
func (h *Handler) DeletePerson(w http.ResponseWriter, r *http.Request) {
	id := parseInt(mux.Vars(r)["id"], -1)
	if id < 0 {
		respondError(w, http.StatusBadRequest, "Invalid ID")
		return
	}
	logger.Log.Infow("DeletePerson", "id", id)

	if err := h.Repo.Delete(id); err != nil {
		logger.Log.Errorw("DB delete failed", "error", err)
		respondError(w, http.StatusNotFound, "Person not found")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func enrichPerson(p *model.Person) error {
	if age, err := service.GetAge(p.Name); err == nil {
		p.Age = age
	} else {
		return err
	}
	if g, err := service.GetGender(p.Name); err == nil {
		p.Gender = g
	}
	if n, err := service.GetNationality(p.Name); err == nil {
		p.Nationality = n
	}
	return nil
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
