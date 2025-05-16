package repository

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"
)

func TestUpdateName(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer func() { _ = db.Close() }()
	repo := NewPeopleRepo(sqlx.NewDb(db, "sqlmock"))

	mock.ExpectQuery(`UPDATE people`).
		WithArgs("Ivan", "Petrov", "", "uuid-1").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow("uuid-1"))

	ok, err := repo.UpdateName(context.Background(), "uuid-1", "Ivan", "Petrov", "")
	require.NoError(t, err)
	require.True(t, ok)
}
