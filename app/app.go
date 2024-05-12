package app

import (
	"clean-architecture/config"
	"clean-architecture/model/dto"
	"clean-architecture/router"
	"database/sql"
	"errors"
	"flag"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/logger"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func initEnv() (dto.ConfigData, error) {
	var configData dto.ConfigData
	if err := godotenv.Load(".env"); err != nil {
		return configData, err
	}

	if port := os.Getenv("PORT"); port != "" {
		configData.AppConfig.Port = port
	}

	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	dbName := os.Getenv("DB_NAME")
	dbMaxIdle := os.Getenv("MAX_IDLE")
	dbMaxConn := os.Getenv("MAX_CONN")
	dbMaxLifeTime := os.Getenv("MAX_LIFE_TIME")
	logMode := os.Getenv("LOG_MODE")

	if dbHost == "" || dbPort == "" || dbUser == "" || dbPass == "" || dbName == "" || dbMaxIdle == "" || dbMaxConn == "" || dbMaxLifeTime == "" || logMode == "" {
		return configData, errors.New("DB Config is not set")
	}

	var err error
	configData.DbConfig.MaxConn, err = strconv.Atoi(dbMaxConn)
	if err != nil {
		return configData, err
	}

	configData.DbConfig.MaxIdle, err = strconv.Atoi(dbMaxIdle)
	if err != nil {
		return configData, err
	}

	configData.DbConfig.Host = dbHost
	configData.DbConfig.DbPort = dbPort
	configData.DbConfig.User = dbUser
	configData.DbConfig.Pass = dbPass
	configData.DbConfig.Database = dbName
	configData.DbConfig.MaxLifeTime = dbMaxLifeTime
	configData.DbConfig.LogMode, err = strconv.Atoi(logMode)
	if err != nil {
		return configData, err
	}
	return configData, nil
}

func initializeDomainModule(r *gin.Engine, db *sql.DB) {
	apiGroup := r.Group("/api")
	v1Group := apiGroup.Group("/v1")
	router.InitRoute(v1Group, db)
}

func RunService() {

	// load env
	configData, err := initEnv()
	if err != nil {
		log.Error().Msg(err.Error())
		return
	}

	conn, err := config.ConnectToDB(configData, log.Logger)
	if err != nil {
		log.Error().Msg("RunService.ConnectToDB.err : " + err.Error())
	}

	duration, err := time.ParseDuration(configData.DbConfig.MaxLifeTime)
	if err != nil {
		log.Error().Msg("RunService.duration.err : " + err.Error())
	}

	// set max conn
	conn.SetConnMaxLifetime(duration)
	conn.SetMaxIdleConns(configData.DbConfig.MaxIdle)
	conn.SetMaxOpenConns(configData.DbConfig.MaxConn)

	defer func() {
		errClose := conn.Close()
		if errClose != nil {
			log.Error().Msg(errClose.Error())
		}
	}()

	time.Local = time.FixedZone("Asia/Jakarta", 7*60*59)
	r := gin.New()
	r.Use(cors.New(cors.Config{
		AllowAllOrigins: false,
		AllowOrigins:    []string{"*"},
		AllowMethods:    []string{"POST", "DELETE", "GET", "OPTIONS", "PUT"},
		AllowHeaders: []string{
			"Origin", "Content-Type",
			"Authorization",
		},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           120 * time.Second,
	}))

	log.Logger = log.With().Caller().Logger()

	r.Use(logger.SetLogger(
		logger.WithLogger(func(_ *gin.Context, l zerolog.Logger) zerolog.Logger {
			return l.Output(os.Stdout).With().Logger()
		}),
	))

	r.Use(gin.Recovery())

	initializeDomainModule(r, conn)

	version := "0.0.1"
	log.Info().Msg(fmt.Sprintf("Service Running version %s", version))
	addr := flag.String("port: ", ":"+configData.AppConfig.Port, "Address to listen and serve")
	if err := r.Run(*addr); err != nil {
		log.Error().Msg(err.Error())
		return
	}
}
