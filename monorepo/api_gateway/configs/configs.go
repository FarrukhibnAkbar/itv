package configs

import (
	helper "itv/monorepo/library/helper"
	"errors"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/joho/godotenv"
	"github.com/spf13/cast"
)

var (
	conf *Configuration
	once sync.Once
)

// Config loads configuration using atomic pattern
func Config() *Configuration {
	once.Do(func() {
		conf = load()
		err := conf.validate()
		if err != nil {
			helper.SendInfo(helper.TgErrorBody{
				Gateway: "in Config file",
				Source:  "api_gateway",
				ErrText: err.Error(),
				Time:    time.Now().Format(time.RFC3339),
			})
			log.Fatalf("error loading configuration: %v", err)
		}
	})
	return conf
}

func load() *Configuration {
	c := Configuration{}
	err := godotenv.Load("./monorepo/api_gateway/.env")
	if err != nil {
		helper.SendInfo(helper.TgErrorBody{
			Gateway: "in Config file",
			Source:  "api_gateway",
			ErrText: err.Error(),
			Time:    time.Now().Format(time.RFC3339),
		})
		fmt.Println("No .env file found")
	}
	c.Environment = cast.ToString(getOrReturnDefault("ENVIRONMENT", "dev"))
	c.LogLevel = cast.ToString(getOrReturnDefault("LOG_LEVEL", "debug"))
	c.MovieServiceHost = cast.ToString(getOrReturnDefault("MOVIE_SERVICE_HOST", "localhost"))
	c.MovieServicePort = cast.ToInt(getOrReturnDefault("MOVIE_SERVICE_PORT", 8090))
	c.HTTPPort = cast.ToString(getOrReturnDefault("HTTP_PORT", ":8080"))
	c.JWTSecretKey = cast.ToString(getOrReturnDefault("JWT_SECRET_KEY", "hellonewworldsecret_key"))

	return &c
}

// Configuration ...
type Configuration struct {
	HTTPPort    string `json:"http_port"`
	LogLevel    string `json:"log_level"`
	Environment string `json:"environment"`

	ServerPort                 int
	ServerHost                 string
	ServiceDir                 string
	AccessTokenDuration        time.Duration
	RefreshTokenDuration       time.Duration
	RefreshPasswdTokenDuration time.Duration
	MovieServiceHost      string
	MovieServicePort      int
	JWTSecretKey          string
}

func (c *Configuration) validate() error {
	if c.HTTPPort == "" {
		return errors.New("http_port required")
	}
	if c.Environment == "" {
		return errors.New("ENVIRONMENT variable required")
	}


	return nil
}

func getOrReturnDefault(key string, defaultValue interface{}) interface{} {
	if os.Getenv(key) == "" {
		return defaultValue
	}
	return os.Getenv(key)
}
