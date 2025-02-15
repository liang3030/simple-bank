package api

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	db "github.com/liang3030/simple-bank/db/sqlc"
	"github.com/liang3030/simple-bank/token"
	"github.com/liang3030/simple-bank/util"
)

// Server HTTP requests for banking service.
type Server struct {
	config     util.Config
	store      db.IStore
	router     *gin.Engine
	tokenMaker token.Maker
}

func NewServer(config util.Config, store db.IStore) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, err
	}
	server := &Server{store: store, config: config, tokenMaker: tokenMaker}

	// Register custom validation functions
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validCurrency)
	}

	server.setupRouter()
	return server, nil
}

func (server *Server) setupRouter() {
	router := gin.Default()
	router.POST("/users", server.createUser)
	router.POST("/users/login", server.loginUser)
	router.POST("/users/renew_access", server.renewAccessToken)

	authRoutes := router.Group("/").Use(authMiddleware(server.tokenMaker))

	authRoutes.GET("/users/:username", server.getUser)
	authRoutes.POST("/accounts", server.createAccount)
	authRoutes.GET("/accounts/:id", server.getAccount)
	authRoutes.GET("/accounts", server.listAccounts)
	authRoutes.PATCH("/accounts/:id", server.updateAccount)

	authRoutes.POST("/transfer", server.transfer)

	server.router = router
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}
