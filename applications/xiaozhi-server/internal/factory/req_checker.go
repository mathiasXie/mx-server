package factory

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func RequestCheck(ctx *gin.Context, req interface{}) string {
	err := ctx.ShouldBindJSON(req)

	if err != nil {
		var validationErrors validator.ValidationErrors
		if errors.As(err, &validationErrors) {
			return validationErrors.Error()
		}
		return err.Error()
	}
	return ""
}
