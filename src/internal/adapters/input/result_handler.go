package input

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"ports-adapters-study/src/internal/core/application"
	dto "ports-adapters-study/src/internal/core/application/dto"
)

type ResultHandler struct {
	resultService application.ResultService
}

func NewResultHandler(service application.ResultService) *ResultHandler {
	return &ResultHandler{
		service,
	}
}

func (r *ResultHandler) GetResult(ctx *gin.Context) error {
	result, err := r.resultService.GetResult()
	if err != nil {
		return err
	}
	ctx.JSON(http.StatusOK, dto.ToResultDto(*result))
	return nil
}

func (r *ResultHandler) GetResultsHistory(ctx *gin.Context) error {
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
