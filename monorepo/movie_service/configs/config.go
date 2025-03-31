package configs

import (
	"errors"
	"fmt"
	"itv/monorepo/library/helper"
	"log"
	"os"
	"sync"
	"time"

	"github.com/joho/godotenv"

	"github.com/spf13/cast"
)

var (
	instance *Configuration
	once     sync.Once
)

// Config ...
func Config() *Configuration {
	fmt.Println("Config")

	once.Do(func() {
		instance = load()
		err := instance.Validate()
		if err != nil {
			helper.SendInfo(helper.TgErrorBody{
				Gateway: "in Config file",
				Source:  "movie_service",
				ErrText: err.Error(),
				Time:    time.Now().Format(time.RFC3339),
			})
			log.Fatalf("error loading configuration: %v", err)
		}
	})

	return instance
}

// Configuration ...
type Configuration struct {
	LogLevel    string `json:"log_level"`
	Environment string `json:"environment"`

	PostgresHost     string
	PostgresPort     int
	PostgresDatabase string
	PostgresUser     string
	PostgresPassword string
	ServiceHost      string
	ServicePort      int
	WebSocketPort    string
	WebSocketHost    string

	RPCPort string

	// context timeout in seconds
	ServerReadTimeout int
	ServiceName       string
}

func load() *Configuration {
	c := Configuration{}
	err := godotenv.Load("./monorepo/movie_service/.env")
	if err != nil {
		helper.SendInfo(helper.TgErrorBody{
			Gateway: "in Config file",
			Source:  "movie_service",
			ErrText: err.Error(),
			Time:    time.Now().Format(time.RFC3339),
		})
		fmt.Println("No .env file found")
	}
	c.Environment = cast.ToString(getOrReturnDefault("ENVIRONMENT", "dev"))
	c.LogLevel = cast.ToString(getOrReturnDefault("LOG_LEVEL", "debug"))
	c.PostgresDatabase = cast.ToString(getOrReturnDefault("POSTGRES_DB", "postgres"))
	c.PostgresUser = cast.ToString(getOrReturnDefault("POSTGRES_USER", "postgres"))
	c.PostgresPassword = cast.ToString(getOrReturnDefault("POSTGRES_PASSWORD", ""))
	c.PostgresHost = cast.ToString(getOrReturnDefault("POSTGRES_HOST", "localhost"))
	c.PostgresPort = cast.ToInt(getOrReturnDefault("POSTGRES_PORT", 5432))
	c.RPCPort = cast.ToString(getOrReturnDefault("RPC_PORT", ":8090"))
	c.ServiceName = cast.ToString(getOrReturnDefault("SERVICE_NAME", "movie_service"))
	c.ServiceHost = cast.ToString(getOrReturnDefault("SERVICE_HOST", "localhost"))
	c.ServicePort = cast.ToInt(getOrReturnDefault("SERVICE_PORT", 8090))

	return &c
}

// Validate validates the configuration
func (c *Configuration) Validate() error {
	if c.RPCPort == "" {
		return errors.New("rpc_port required")
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
