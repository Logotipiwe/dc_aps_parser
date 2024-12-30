package input

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"ports-adapters-study/src/internal/core/application"
	"ports-adapters-study/src/internal/core/application/dto"
	driverport "ports-adapters-study/src/internal/core/ports/input"
	"ports-adapters-study/src/pkg"
)

type resultController struct {
	resultService *application.ResultService
}

func NewResultController(router *gin.Engine, service *application.ResultService) driverport.ResultController {
	r := &resultController{
		service,
	}
	router.GET("/get", pkg.WithError(r.GetResult))
	router.GET("/history", pkg.WithError(r.GetResultsHistory))
	return r
}

func (r *resultController) GetResult(ctx *gin.Context) error {
	result, err := r.resultService.GetResult()
	if err != nil {
		return err
	}
	ctx.JSON(http.StatusOK, dto.ToResultDto(*result))
	return nil
}

func (r *resultController) GetResultsHistory(ctx *gin.Context) error {
	history, err := r.resultService.GetResultHistory()
	if err != nil {
		return err
	}
	historyDtos := make([]dto.ParseResult, len(history))
	for i, item := range history {
		historyDtos[i] = dto.ToResultDto(item)
	}
	ctx.JSON(http.StatusOK, gin.H{
		"results": historyDtos,
	})
	return nil
}
