package main

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"go-template/api"
	"go-template/api/handlers"
	appLogger "go-template/pkg/logger"

	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Server struct {
		RunMode      string        `envconfig:"SERVER_RUN_MODE" default:"debug"`
		Port         int           `envconfig:"SERVER_PORT" default:"8080"`
		ReadTimeout  time.Duration `envconfig:"SERVER_READ_TIMEOUT" default:"600s"`
		WriteTimeout time.Duration `envconfig:"SERVER_WRITE_TIMEOUT" default:"600s"`
	}
	Database struct {
		Host     string `envconfig:"DATABASE_HOST" required:"true"`
		Port     int    `envconfig:"DATABASE_PORT" default:"27017"`
		User     string `envconfig:"DATABASE_USER" required:"true"`
		Password string `envconfig:"DATABASE_PASSWORD" required:"true"`
		Name     string `envconfig:"DATABASE_NAME" default:"admin"`
		TLSMode  bool   `envconfig:"DATABASE_TLS_MODE" default:"true"`
	}
	Logger struct {
		Level  string `envconfig:"LOG_LEVEL" default:"info"`
		Output string `envconfig:"LOG_OUTPUT" default:"console"`
	}
}

const serviceName = "my-service"

func main() {
	// Load the .env file in the current directory
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file", err)
	}

	var cfg Config
	if err := envconfig.Process("", &cfg); err != nil {
		log.Fatalf("start up failed - error parsing app settings:%s", err)
	}

	log.Printf("starting service with config:\n"+
		"SERVER_RUN_MODE: %s\n"+
		"SERVER_PORT: %s\n"+
		"SERVER_READ_TIMEOUT: %s\n"+
		"SERVER_WRITE_TIMEOUT: %s\n"+
		"DATABASE_HOST: %s\n"+
		"DATABASE_PORT: %d\n"+
		"DATABASE_USER: %s\n"+
		"DATABASE_PASSWORD: %s\n"+
		"DATABASE_NAME: %s\n"+
		"DATABASE_TLS_MODE: %t\n"+
		"LOG_LEVEL: %s\n"+
		"LOG_OUTPUT: %s\n",
		cfg.Server.RunMode,
		strconv.Itoa(cfg.Server.Port),
		cfg.Server.ReadTimeout,
		cfg.Server.WriteTimeout,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Name,
		cfg.Database.TLSMode,
		cfg.Logger.Level,
		cfg.Logger.Output,
	)

	gin.SetMode(cfg.Server.RunMode)
	govalidator.SetFieldsRequiredByDefault(true)

	logger, err := appLogger.NewLogger(cfg.Logger, serviceName)
	if err != nil {
		log.Fatalf("failed to create logger: %s", err)
	}

	helloWorldHandler := handlers.NewHelloWorldHandler()

	router := api.NewRouter(logger, helloWorldHandler)
	server := http.Server{
		Addr:         ":" + strconv.Itoa(cfg.Server.Port),
		Handler:      router,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
	}

	if err := server.ListenAndServe(); err != nil {
		logger.Error().Msgf("failed to run server: %s", err)
	}
}
