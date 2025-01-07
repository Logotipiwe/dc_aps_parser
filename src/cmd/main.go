package main

import (
	. "dc-aps-parser/src/internal/adapters/input"
	"dc-aps-parser/src/internal/adapters/output"
	. "dc-aps-parser/src/internal/core/application"
	. "dc-aps-parser/src/internal/core/ports"
	"dc-aps-parser/src/internal/infrastructure"
	"dc-aps-parser/src/internal/infrastructure/tg"
	"github.com/gin-gonic/gin"
	"os"
)

func main() {
	config := infrastructure.NewConfig()
	botAPI := tg.NewBotAPI(os.Getenv("TELEGRAM_TOKEN"))

	adapters := OutputPorts{
		TargetClientPort: output.NewTargetClientWebAdapter(),
		NotificationPort: output.NewNotificationAdapterTg(botAPI),
	}

	resultService := NewResultService(adapters.TargetClientPort)
	adminService := NewAdminService()
	parserNotificationService := NewParserNotificationService(config, adapters.NotificationPort)
	app := App{
		ResultService: resultService,
		ParserService: NewParserService(config, resultService, parserNotificationService),
		AdminService:  adminService,
	}

	router := gin.Default()

	_ = InputAdapters{
		NewParserAdapterHttp(router, app.ParserService),
		NewParserAdapterTg(botAPI, app.ParserService, adminService),
	}

	err := router.Run(":81")
	if err != nil {
		panic(err.Error())
	}
}

type InputAdapters struct {
	*ParserAdapterHttp
	*ParserAdapterTg
}
