package handler

import (
	"encoding/json"
	"github.com/51mans0n/effective-mobile-task/internal/repository"
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
