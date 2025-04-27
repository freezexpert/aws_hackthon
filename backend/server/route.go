package server

import (
	"backend/controller"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
)

func (srv *server) routes() http.Handler {

	// srv.router.Use(gin.Logger())
	// srv.router.Use(gin.Recovery())

	srv.router.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"POST, GET, OPTIONS, PUT, DELETE, UPDATE"},
		AllowHeaders:     []string{"Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,

		MaxAge: 12 * time.Hour,
	}))

	controller := &controller.BaseController{
		Service: srv.service,
	}

	v1 := srv.router.Group("/")
	v1.Use()
	{
		v1.POST("/user_history", controller.GetHistory)
		v1.POST("/", controller.PostHistory)
		v1.POST("/generate_response", controller.GenerateResponse)
		v1.POST("/chat", controller.ProcessChat)
	}
	return srv.router
}
