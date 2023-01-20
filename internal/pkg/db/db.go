package db

import (
	"fmt"

	"github.com/aliyevazam/tages/internal/pkg/config"
	"github.com/jmoiron/sqlx"
)

func ConnectToDB(cfg config.Config) (*sqlx.DB, error) {
	psqlString := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.PostgresHost,
		cfg.PostgresPort,
		cfg.PostgresUser,
		cfg.PostgresPassword,
		cfg.PostgresDatabase,
	)

	connDB, err := sqlx.Connect("postgres", psqlString)
	if err != nil {
		fmt.Println("error while connecting to db", err)
		return nil, err
	}
	return connDB, nil
}
