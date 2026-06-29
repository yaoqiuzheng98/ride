package handler

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"ride/db/table"
)

type createUserRequest struct {
	Phone string `json:"phone" binding:"required"`
}

type userResponse struct {
	Id    int64  `json:"id"`
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
		if errors.Is(err, gorm.ErrDuplicatedKey) || strings.Contains(err.Error(), "Duplicate entry") {
			ctx.JSON(http.StatusConflict, gin.H{"error": "phone already exists"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, userResponse{Id: u.Id, Phone: u.Phone})
}
