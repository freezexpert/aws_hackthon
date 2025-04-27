package server

import (
	"backend/middleware/cors"
	"backend/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

type server struct {
	router  *gin.Engine
	service models.Service
}

// NewServer returns new http.Server.
func NewServer(service models.Service) *http.Server {

	// gin.SetMode(gin.TestMode)
	router := gin.New()

	router.HandleMethodNotAllowed = true

	router.Use(
		gin.Logger(),
		gin.Recovery(),
		cors.Default(),
		cors.CORSMiddleware(),
	)

	srv := &server{
		router:  router,
		service: service,
	}

	return &http.Server{
		Addr:    "0.0.0.0:8888",
		Handler: srv.routes(),
	}
}
