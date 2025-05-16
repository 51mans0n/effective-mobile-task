package repository

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"
)

func TestList(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	sqlxdb := sqlx.NewDb(db, "sqlmock")
	repo := NewPeopleRepo(sqlxdb)

	mock.ExpectQuery(`COUNT\(\*\)`).WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))
	mock.ExpectQuery(`SELECT \*`).WillReturnRows(
		sqlmock.NewRows([]string{"id", "name"}).AddRow("uuid-1", "Bob"),
	)

	out, err := repo.List(context.Background(), ListFilter{Page: 1, Limit: 20})
	require.NoError(t, err)
	require.Equal(t, int64(1), out.Total)
	require.Len(t, out.Records, 1)
}
