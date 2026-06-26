package handler

import (
	"net/http"

	"ride/db/table"

	"github.com/gin-gonic/gin"
)

func Point(ctx *gin.Context) {
	points := table.GetPoints().GetPoints()
	ctx.JSON(http.StatusOK, points)
}
