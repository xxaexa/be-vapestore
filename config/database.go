package config

import (
	"clean-architecture/model/dto"
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	"github.com/rs/zerolog"
)

func ConnectToDB(in dto.ConfigData, logger zerolog.Logger) (*sql.DB, error) {
	logger.Info().Msg("Trying connect to db..")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", in.DbConfig.Host, in.DbConfig.User, in.DbConfig.Pass, in.DbConfig.Database, in.DbConfig.DbPort)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to ping db")
		return nil, err
	}

	logger.Info().Msg("Successfully connected to db")
	return db, nil
}
