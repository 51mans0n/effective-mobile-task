package repository

import (
	"context"
	"regexp"
	"testing"

	"github.com/51mans0n/effective-mobile-task/internal/model"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"
)

func TestCreate(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	sqlxdb := sqlx.NewDb(db, "sqlmock")
	repo := NewPeopleRepo(sqlxdb)

	p := &model.Person{Name: "Ivan", Surname: "Petrov"}

	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO people`)).
		WithArgs("Ivan", "Petrov", nil, nil, nil, nil, nil).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow("uuid-123"))

	err := repo.Create(context.Background(), p)
	require.NoError(t, err)
	require.Equal(t, "uuid-123", p.ID)
	require.NoError(t, mock.ExpectationsWereMet())
}
