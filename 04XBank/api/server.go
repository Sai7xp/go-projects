/*
* Created on 06 May 2024
* @author Sai Sumanth
 */
package api

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	db "github.com/sai7xp/xbank/db/sqlc"
)

// Server serves HTTP requests for out banking service.
type Server struct {
	store  *db.Store
	router *gin.Engine
}

// NewServer creates a new HTTP server and setup routing.
func NewServer(store *db.Store) *Server {
	server := &Server{store: store}
	router := gin.Default()

	if validate, ok := binding.Validator.Engine().(*validator.Validate); ok {
		validate.RegisterValidation("currency", validCurrency)
	}

	// Account API Routes
	router.POST("/accounts", server.createAccount)
	router.GET("/accounts/:id", server.getAccount)
	router.GET("/accounts", server.listAccounts)
	router.DELETE("/accounts/:id", server.deleteAccount)
	router.PATCH("/accounts/:id", server.updateAccount)

	server.router = router
	return server
}

// Start - runs the HTTP server on a specific address
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	if err == nil {
		return gin.H{"error": "Null pointer execption"}
	}
	return gin.H{"error": err.Error()}
}
