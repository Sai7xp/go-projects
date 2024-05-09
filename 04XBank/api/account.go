/*
* Created on 06 May 2024
* @author Sai Sumanth
 */
package api

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	db "github.com/sai7xp/xbank/db/sqlc"
)

type createAccountRequest struct {
	Owner    string `json:"owner" binding:"required"`
	Currency string `json:"currency" binding:"required,currency"` // "currency" is a custom validator(registered in server.go)
}

// createAccount API Handler function
func (server *Server) createAccount(ctx *gin.Context) {
	var reqData createAccountRequest
	if err := ctx.ShouldBindJSON(&reqData); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	params := db.CreateAccountParams{
		Owner:    reqData.Owner,
		Balance:  0,
		Currency: reqData.Currency,
	}

	account, err := server.store.CreateAccount(ctx, params)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
	}
	ctx.JSON(http.StatusOK, account)
}

type getAccountRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

// getAccount API Handler function
func (server *Server) getAccount(ctx *gin.Context) {
	var getAccReq getAccountRequest
	if err := ctx.BindUri(&getAccReq); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	account, err := server.store.GetAccount(ctx, getAccReq.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, map[string]interface{}{"success": true, "data": account})
}

type listAccountsRequest struct {
	PageId   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=10"`
}

// listAccounts API Handler Function
func (server *Server) listAccounts(ctx *gin.Context) {
	var req listAccountsRequest
	if err := ctx.BindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	params := db.ListAccountsParams{
		Offset: (req.PageId - 1) * req.PageSize,
		Limit:  req.PageSize,
	}
	accounts, err := server.store.ListAccounts(ctx, params)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, map[string]interface{}{"success": true, "data": accounts})
}

type deleteAccountRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

// deleteAccount API Handler Function
func (server *Server) deleteAccount(ctx *gin.Context) {
	var req deleteAccountRequest
	if err := ctx.BindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	rowsAffected, err := server.store.DeleteAccount(ctx, req.ID)
	if rowsAffected == 0 {
		errMessage := fmt.Sprintf("Account with Id : %d Not found in the db.", req.ID)
		ctx.JSON(http.StatusNotFound, map[string]string{"message": errMessage})
		return
	}
	if err != nil {
		fmt.Println(err)
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, map[string]interface{}{"message": "Account Deleted Successfully!"})
}

// updateAccount API Handler Function
func (server *Server) updateAccount(ctx *gin.Context) {
	// TODO: Invoke UpdateAccount
}
