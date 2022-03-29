package domain

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/prosperitybot/worker/internal"

	_ "github.com/go-sql-driver/mysql"
)

type TestRepositoryDatabase struct {
	db *sqlx.DB
}

func (d TestRepositoryDatabase) Test(ctx context.Context) ([]string, []string) {
	return nil, nil
}

func NewTestRepositoryDatabase() TestRepositoryDatabase {
	return TestRepositoryDatabase{internal.Database}
}
