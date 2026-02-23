package app

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"

	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
	"github.com/go-sql-driver/mysql"
	"offgrocery-assessment/internal/config"
	"offgrocery-assessment/internal/ent"
)

func NewDB(_ context.Context, cfg config.Config, configFuncs ...func(c *mysql.Config)) (*sql.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName,
	)

	db, err := sql.Open(dialect.MySQL, BuildDSN(dsn, configFuncs...))
	if err != nil {
		slog.Error("components: failed to open database connection", "error", err)
		return nil, err
	}

	return db, nil
}

func NewEntClient(db *sql.DB) *ent.Client {
	dbDriver := entsql.OpenDB(dialect.MySQL, db)
	return ent.NewClient(ent.Driver(dbDriver))
}

func BuildDSN(dbURL string, configFuncs ...func(c *mysql.Config)) string {
	cfg, err := mysql.ParseDSN(dbURL)
	if err != nil {
		slog.Error("components: failed to parse DSN", "error", err)
		return dbURL
	}

	if cfg.Params == nil {
		cfg.Params = make(map[string]string)
	}

	for _, configure := range configFuncs {
		configure(cfg)
	}

	return cfg.FormatDSN()
}

func ConfigureMySQLForMigration(c *mysql.Config) {
	c.MultiStatements = true
}

func ConfigureMySQLParseTime(c *mysql.Config) {
	c.ParseTime = true
}
