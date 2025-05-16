// Package repository implements PostgreSQL access using sqlx and squirrel.
package repository

import (
	"context"
	"database/sql"
	"errors"
	"github.com/51mans0n/effective-mobile-task/internal/model"

	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
)

// PeopleRepo defines data-access interface for people table.
type PeopleRepo interface {
	Create(ctx context.Context, p *model.Person) error
	GetByID(ctx context.Context, id string) (*model.Person, error)
	List(ctx context.Context, f ListFilter) (*PaginatedPeople, error)
	UpdateName(ctx context.Context, id string, name, surname, patronymic string) (bool, error)
	Delete(ctx context.Context, id string) (bool, error)
}

type repo struct {
	db *sqlx.DB
	sb squirrel.StatementBuilderType
}

// ListFilter defines filters and pagination parameters for listing people.
type ListFilter struct {
	Name    string
	Country string
	Gender  string
	Page    int // 1-based
	Limit   int
}

// PaginatedPeople represents a paginated response with total count.
type PaginatedPeople struct {
	Total   int64
	Records []*model.Person
}

// NewPeopleRepo returns a new repository backed by sqlx and squirrel.
func NewPeopleRepo(db *sqlx.DB) PeopleRepo {
	return &repo{
		db: db,
		sb: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}
}

func (r *repo) Create(ctx context.Context, p *model.Person) error {
	query, args, err := r.sb.
		Insert("people").
		Columns("name", "surname", "patronymic", "age", "gender", "country_code", "nat_probability").
		Values(p.Name, p.Surname, p.Patronymic, p.Age, p.Gender, p.CountryCode, p.NatProbability).
		Suffix("RETURNING id").
		ToSql()
	if err != nil {
		return err
	}
	return r.db.QueryRowContext(ctx, query, args...).Scan(&p.ID)
}

func (r *repo) GetByID(ctx context.Context, id string) (*model.Person, error) {
	query, args, err := r.sb.
		Select("*").
		From("people").
		Where(squirrel.Eq{"id": id}).
		ToSql()
	if err != nil {
		return nil, err
	}
	var p model.Person
	if err := r.db.GetContext(ctx, &p, query, args...); err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *repo) List(ctx context.Context, f ListFilter) (*PaginatedPeople, error) {
	sb := r.sb.Select("*").From("people")
	if f.Name != "" {
		sb = sb.Where("name ILIKE ?", "%"+f.Name+"%")
	}
	if f.Country != "" {
		sb = sb.Where(squirrel.Eq{"country_code": f.Country})
	}
	if f.Gender != "" {
		sb = sb.Where(squirrel.Eq{"gender": f.Gender})
	}

	// подсчёт total
	countSb := r.sb.Select("COUNT(*)").From("people")
	if f.Name != "" {
		countSb = countSb.Where("name ILIKE ?", "%"+f.Name+"%")
	}
	if f.Country != "" {
		countSb = countSb.Where(squirrel.Eq{"country_code": f.Country})
	}
	if f.Gender != "" {
		countSb = countSb.Where(squirrel.Eq{"gender": f.Gender})
	}

	countQuery, countArgs, _ := countSb.ToSql()
	var total int64
	if err := r.db.GetContext(ctx, &total, countQuery, countArgs...); err != nil {
		return nil, err
	}

	// пагинация
	offset := (f.Page - 1) * f.Limit
	sb = sb.Limit(uint64(f.Limit)).Offset(uint64(offset)).OrderBy("created_at DESC")

	query, args, _ := sb.ToSql()
	var list []*model.Person
	if err := r.db.SelectContext(ctx, &list, query, args...); err != nil {
		return nil, err
	}

	return &PaginatedPeople{Total: total, Records: list}, nil
}

func (r *repo) UpdateName(ctx context.Context, id, name, surname, patr string) (bool, error) {
	qb := r.sb.
		Update("people").
		Set("name", name).
		Set("surname", surname).
		Set("patronymic", patr).
		Set("updated_at", squirrel.Expr("now()")).
		Where(squirrel.Eq{"id": id}).
		Suffix("RETURNING id")

	q, args, _ := qb.ToSql()
	var dummy string
	err := r.db.QueryRowContext(ctx, q, args...).Scan(&dummy)
	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	}
	return err == nil, err
}

func (r *repo) Delete(ctx context.Context, id string) (bool, error) {
	q, args, _ := r.sb.
		Delete("people").
		Where(squirrel.Eq{"id": id}).
		Suffix("RETURNING id").
		ToSql()
	var dummy string
	err := r.db.QueryRowContext(ctx, q, args...).Scan(&dummy)
	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	}
	return err == nil, err
}
