package api

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	db "github.com/liang3030/simple-bank/db/sqlc"
)

// Server HTTP requests for banking service.
type Server struct {
	store  db.IStore
	router *gin.Engine
}

func NewServer(store db.IStore) *Server {
	server := &Server{store: store}
	router := gin.Default()

	// Register custom validation functions
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validCurrency)
	}
	router.POST("/users", server.createUser)
	router.GET("/users/:username", server.getUser)
	router.POST("/accounts", server.createAccount)
	router.GET("/accounts/:id", server.getAccount)
	router.GET("/accounts", server.listAccounts)
	router.PATCH("/accounts/:id", server.updateAccount)
	router.POST("/transfer", server.transfer)

	server.router = router
	return server
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}
