package input

import (
	"github.com/gin-gonic/gin"
	"ports-adapters-study/src/internal/core/application"
	driverport "ports-adapters-study/src/internal/core/ports/input"
	"ports-adapters-study/src/pkg"
	"strconv"
)

type parserController struct {
	*application.ParserService
}

func NewParserController(
	router *gin.Engine,
	service *application.ParserService,
) driverport.ParserController {
	p := &parserController{
		service,
	}

	parserApi := router.Group("/parser")
	parserApi.POST("/new", pkg.WithError(p.NewParser))
	parserApi.POST("/stop", pkg.WithError(p.StopParser))

	return p
}

func (c *parserController) NewParser(ctx *gin.Context) error {
	parser, err := c.ParserService.NewParser()
	if err != nil {
		return err
	}
	ctx.JSON(200, parser)
	return nil
}

func (c *parserController) StopParser(ctx *gin.Context) error {
	idStr := ctx.Query("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return err
	}
	err = c.ParserService.StopParser(id)
	if err != nil {
		return err
	}
	ctx.JSON(200, nil)
	return nil
}
