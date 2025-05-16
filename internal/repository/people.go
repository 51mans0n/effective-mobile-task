package repository

import (
	"context"
	"github.com/51mans0n/effective-mobile-task/internal/model"

	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
)

type PeopleRepo interface {
	Create(ctx context.Context, p *model.Person) error
	GetByID(ctx context.Context, id string) (*model.Person, error)
	List(ctx context.Context, f ListFilter) (*PaginatedPeople, error)
}

type repo struct {
	db *sqlx.DB
	sb squirrel.StatementBuilderType
}

type ListFilter struct {
	Name    string
	Country string
	Gender  string
	Page    int // 1-based
	Limit   int
}

type PaginatedPeople struct {
	Total   int64
	Records []*model.Person
}

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
