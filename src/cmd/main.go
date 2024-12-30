package main

import (
	"github.com/gin-gonic/gin"
	"ports-adapters-study/src/internal/adapters/input"
	"ports-adapters-study/src/internal/adapters/output"
	"ports-adapters-study/src/internal/core/application"
	driverport "ports-adapters-study/src/internal/core/ports/input"
	krishawebclient "ports-adapters-study/src/internal/infrastructure/krisha"
	"ports-adapters-study/src/internal/infrastructure/tg"
)

type services struct {
	driverport.ResultController
	driverport.ParserController
}

func initServices(router *gin.Engine) services {
	resultService := application.NewResultService(
		output.NewResultRepository(),
		output.NewKrishaWebClientAdapter(
			krishawebclient.NewKrishaWebClient(),
		),
	)
	s := services{
		input.NewResultController(
			router,
			resultService,
		),
		input.NewParserController(
			router,
			application.NewParserService(
				resultService, output.NewTgClientAdapter(
					tg.NewTgClient(),
				),
			),
		),
	}
	return s
}

func main() {
	println("Started!")

	router := gin.Default()

	initServices(router)

	err := router.Run(":81")
	if err != nil {
		panic(err.Error())
	}
}
