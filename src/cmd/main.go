package main

import (
	"github.com/gin-gonic/gin"
	"ports-adapters-study/src/internal/adapters/input"
	"ports-adapters-study/src/internal/adapters/output"
	"ports-adapters-study/src/internal/core/application"
	driverport "ports-adapters-study/src/internal/core/ports/input"
	drivenport "ports-adapters-study/src/internal/core/ports/output"
)

type Controllers struct {
	router *gin.Engine
	driverport.ResultController
	driverport.ParserController
}

type App struct {
	*application.ResultService
	*application.ParserService
}

type Adapters struct {
	drivenport.ResultDB
	drivenport.TargetClient
	drivenport.NotificationClient
}

func CreateProdAdapters() Adapters {
	return Adapters{
		ResultDB:           output.NewResultRepository(),
		TargetClient:       output.NewKrishaWebClientAdapter(),
		NotificationClient: output.NewTgClientAdapter(),
	}
}

func NewApp(
	a Adapters,
) App {
	resultService := application.NewResultService(a.ResultDB, a.TargetClient)

	return App{
		ResultService: resultService,
		ParserService: application.NewParserService(resultService, a.NotificationClient),
	}
}

func InitProdControllers(app App) *gin.Engine {
	router := gin.Default()
	_ = Controllers{
		router,
		input.NewResultController(router, app.ResultService),
		input.NewParserController(router, app.ParserService),
	}
	return router
}

func main() {
	println("Started!")

	app := NewApp(CreateProdAdapters())

	router := InitProdControllers(app)

	err := router.Run(":81")
	if err != nil {
		panic(err.Error())
	}
}
