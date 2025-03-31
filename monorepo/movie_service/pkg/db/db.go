package db

import (
	"itv/monorepo/movie_service/configs"
	"itv/monorepo/library/helper"
	"fmt"
	"log"
	"os"
	"time"

	// needed migration packages
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/joho/godotenv/autoload" // load .env file automatically
	"go.uber.org/zap"

	// database driver
	_ "github.com/lib/pq"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Init initializes database connection then connect with postgres
func Init(config *configs.Configuration) (*gorm.DB, error) {

	fmt.Println("init db")
	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", config.PostgresUser, config.PostgresPassword, config.PostgresHost, config.PostgresPort, config.PostgresDatabase)

	m, err := migrate.New(
		"file://monorepo/movie_service/pkg/db/migrations", 
		dbURL,
	)
	if err != nil {
		helper.SendInfo(helper.TgErrorBody{
			Gateway: "in movie_service/pkg/db file",
			Source:  "movie_service",
			ErrText: err.Error(),
			Time:    time.Now().Format(time.RFC3339),
		})
		log.Fatal("error in creating migrations: ", zap.Error(err))
	}
	fmt.Printf("")
	if err = m.Up(); err != nil {
		helper.SendInfo(helper.TgErrorBody{
			Gateway: "in movie_service/pkg/db file",
			Source:  "movie_service",
			ErrText: err.Error(),
			Time:    time.Now().Format(time.RFC3339),
		})
		log.Println("error updating migrations: ", zap.Error(err))
	}

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second,  // Slow SQL threshold
			LogLevel:                  logger.Error, // For debugging you can set Info level
			IgnoreRecordNotFoundError: true,         // Ignore ErrRecordNotFound error for logger
			ParameterizedQueries:      false,        // Don't include params in the SQL log
			Colorful:                  true,         // Disable color
		},
	)

	db, err := gorm.Open(postgres.Open(dbURL), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		helper.SendInfo(helper.TgErrorBody{
			Gateway: "in movie_service/pkg/db file",
			Source:  "movie_service",
			ErrText: err.Error(),
			Time:    time.Now().Format(time.RFC3339),
		})
		fmt.Println(err.Error())
		panic("Failed to connect to database")
	}
	sqlDB, err := db.DB()
	if err != nil {
		helper.SendInfo(helper.TgErrorBody{
			Gateway: "in movie_service/pkg/db file",
			Source:  "movie_service",
			ErrText: err.Error(),
			Time:    time.Now().Format(time.RFC3339),
		})
		log.Println("Failed to get to database", err)
	} else {
		sqlDB.SetMaxOpenConns(100)
		sqlDB.SetMaxIdleConns(100)
		sqlDB.SetConnMaxLifetime(time.Minute * 5)
	}

	return db, nil
}
