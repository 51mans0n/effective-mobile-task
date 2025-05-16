package handler

import (
	"encoding/json"
	"github.com/51mans0n/effective-mobile-task/internal/repository"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/51mans0n/effective-mobile-task/internal/model"
	"github.com/51mans0n/effective-mobile-task/internal/service"
	"github.com/go-playground/validator/v10"
)

type createReq struct {
	Name     string `json:"name"      validate:"required,alpha"`
	Surname  string `json:"surname"   validate:"required,alpha"`
	Patronym string `json:"patronymic,omitempty"`
}

type updateReq struct {
	Name     string `json:"name"      validate:"required,alpha"`
	Surname  string `json:"surname"   validate:"required,alpha"`
	Patronym string `json:"patronymic"`
}

type listMeta struct {
	Page  int   `json:"page"`
	Limit int   `json:"limit"`
	Total int64 `json:"total"`
}
type listResp struct {
	Meta listMeta        `json:"meta"`
	Data []*model.Person `json:"data"`
}

var v = validator.New() // можно передать извне

// CreatePerson godoc
// @Summary Create a new person
// @Description Creates a new person and enriches data from external APIs (agify, genderize, nationalize)
// @Tags People
// @Accept json
// @Produce json
// @Param body body createReq true "Person data to create"
// @Success 201 {object} model.Person
// @Failure 400 {string} string "Bad request"
// @Failure 500 {string} string "Internal error"
// @Router /people [post]
func Create(svc *service.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req createReq
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "bad json: "+err.Error(), http.StatusBadRequest)
			return
		}
		if err := v.Struct(req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		p := &model.Person{
			Name:    req.Name,
			Surname: req.Surname,
		}
		if req.Patronym != "" {
			p.Patronymic = &req.Patronym
		}

		if err := svc.Create(r.Context(), p); err != nil {
			log.Printf("create failed: %v", err)
			http.Error(w, "internal error", http.StatusInternalServerError)

			return
		}

		w.Header().Set("Location", "/people/"+p.ID)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(p)
	}
}

// ListPeople godoc
// @Summary Get all people
// @Description Returns a paginated list of people, with optional filters by name, gender, and country
// @Tags People
// @Param name query string false "Search by name (case-insensitive, partial match)"
// @Param gender query string false "Gender filter (male or female)"
// @Param country query string false "Country code (ISO-2)"
// @Param page query int false "Page number (>=1)" default(1)
// @Param limit query int false "Results per page (1-100)" default(20)
// @Success 200 {object} listResp
// @Failure 500 {string} string "Internal error"
// @Router /people [get]
func List(svc *service.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()

		page, _ := strconv.Atoi(q.Get("page"))
		limit, _ := strconv.Atoi(q.Get("limit"))
		if page <= 0 {
			page = 1
		}
		if limit <= 0 || limit > 100 {
			limit = 20
		}

		f := repository.ListFilter{
			Name:    q.Get("name"),
			Country: strings.ToUpper(q.Get("country")),
			Gender:  q.Get("gender"),
			Page:    page,
			Limit:   limit,
		}
		res, err := svc.List(r.Context(), f)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		resp := listResp{
			Meta: listMeta{Page: page, Limit: limit, Total: res.Total},
			Data: res.Records,
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(resp)
	}
}

// GetByID godoc
// @Summary Get a person by ID
// @Param id path string true "UUID of the person"
// @Success 200 {object} model.Person
// @Failure 404 {string} string "not found"
// @Router /people/{id} [get]
func GetByID(svc *service.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		rec, err := svc.Get(r.Context(), id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if rec == nil {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(rec)
	}
}

// Update godoc
// @Summary Update person's name
// @Param id path string true "UUID"
// @Param body body updateReq true "updated names"
// @Success 204 {string} string "no content"
// @Failure 404 {string} string "not found"
// @Router /people/{id} [put]
func Update(svc *service.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		var req updateReq
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "bad json", http.StatusBadRequest)
			return
		}
		if err := v.Struct(req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		ok, err := svc.UpdateName(r.Context(), id, req.Name, req.Surname, req.Patronym)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if !ok {
			http.NotFound(w, r)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}

// Delete godoc
// @Summary Delete a person
// @Param id path string true "UUID"
// @Success 204 {string} string "deleted"
// @Failure 404 {string} string "not found"
// @Router /people/{id} [delete]
func Delete(svc *service.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		ok, err := svc.Delete(r.Context(), id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if !ok {
			http.NotFound(w, r)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}
