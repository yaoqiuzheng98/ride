package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"ride/db/table"
)

func Point(ctx *gin.Context) {
	points := table.GetPoints().GetPoints()
	ctx.JSON(http.StatusOK, points)
}
