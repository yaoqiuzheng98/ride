package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"ride/db/collection"
)

func Point(ctx *gin.Context) {
	points := collection.GetPoints().GetPoints()
	ctx.JSON(http.StatusOK, points)
}
