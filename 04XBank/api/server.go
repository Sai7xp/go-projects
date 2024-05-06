/*
* Created on 06 May 2024
* @author Sai Sumanth
 */
package api

import (
	"github.com/gin-gonic/gin"
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

	router.POST("/accounts", server.createAccount)
	router.GET("/account/:id", server.getAccount)
	router.GET("/accounts", server.listAccounts)
	router.DELETE("/account/:id", server.deleteAccount)

	server.router = router
	return server
}

// Start - runs the HTTP server on a specific address
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
