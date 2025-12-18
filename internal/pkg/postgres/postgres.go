package postgres

import (
	"database/sql"
	"main/internal/pkg/config"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/extra/bundebug"
)

func NewDB() *bun.DB {

	dsn := "postgres://" + config.GetConfig().DBUsername + ":" + config.GetConfig().DBPassword + "@localhost:5432/" + config.GetConfig().DBName + "?sslmode=disable"

	sqlDb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))

	db := bun.NewDB(sqlDb, pgdialect.New())

	db.AddQueryHook(bundebug.NewQueryHook(
		bundebug.WithVerbose(true),
		bundebug.FromEnv("BUNDEBUG"),
	))

	return db
}
