package api

import (
	"go-template/api/handlers"
	"go-template/api/middleware"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

func NewRouter(logger zerolog.Logger, helloWorldHandler *handlers.HelloWorldHandler) *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(middleware.Logging(logger, nil))

	helloWorldGroup := r.Group("/hello-world")
	{
		helloWorldGroup.GET("", helloWorldHandler.GetHelloWorld)
	}

	return r
}
