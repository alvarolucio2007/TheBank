package api

import (
	db "github.com/alvarolucio2007/TheBank/db/sqlc"
	"github.com/gin-gonic/gin"
)

type Server struct {
	store  db.Store
	router *gin.Engine
}

func NewServer(store db.Store) (*Server, error) {
	server := &Server{store: store}
	router := gin.Default()
	server.router = router
	server.setupRouter()
	return server, nil
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func (server *Server) setupRouter() {
	router := gin.Default()

	router.POST("/users", server.createUser)

	router.POST("/accounts", server.createAccount)
	router.GET("/accounts/:id", server.getAccount)
	router.GET("/accounts", server.listAccount)

	router.POST("/transfers", server.createTransfer)

	server.router = router
}

func errorResponse(err error) gin.H {
	if err == nil {
		return gin.H{"error": "unknown error"}
	}
	return gin.H{"error": err.Error()}
}
