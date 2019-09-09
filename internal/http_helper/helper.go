// Author: xufei
// Date: 2019-09-06 17:02

package http_helper

import (
	"gim/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Render400(ctx *gin.Context, err error) {
	ctx.JSON(http.StatusBadRequest, model.ErrResult{Message: err.Error()})
}

func Render500(ctx *gin.Context, err error) {
	ctx.JSON(http.StatusInternalServerError, model.ErrResult{Message: err.Error()})
}

func RenderOK(ctx *gin.Context, data interface{}) {
	ctx.JSON(http.StatusOK, model.CodeResult{
		Code:    "0",
		Message: "OK",
		Data:    data,
	})
}

func RenderCreated(ctx *gin.Context, data interface{}) {
	ctx.JSON(http.StatusCreated, model.CodeResult{
		Code:    "0",
		Message: "OK",
		Data:    data,
	})
}
