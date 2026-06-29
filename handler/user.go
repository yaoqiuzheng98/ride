package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"ride/db/table"
)

type createUserRequest struct {
	Phone string `json:"phone" binding:"required"`
}

type userResponse struct {
	Id    uint   `json:"id"`
	BizID string `json:"biz_id"`
	Phone string `json:"phone"`
}

// CreateUser 输入手机号创建一个用户。
// POST /user  body: {"phone":"..."}
func CreateUser(ctx *gin.Context) {
	var req createUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "phone required"})
		return
	}
	u, err := table.CreateUser(req.Phone)
	if err != nil {
		if errors.Is(err, table.ErrPhoneAlreadyExists) {
			ctx.JSON(http.StatusConflict, gin.H{"error": "phone already exists"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, userResponse{Id: u.ID, BizID: u.BizID, Phone: u.Phone})
}
