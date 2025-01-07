package tests

import (
	"dc-aps-parser/src/internal/adapters/output"
	"dc-aps-parser/src/internal/core/application"
	"dc-aps-parser/src/internal/core/domain"
	"dc-aps-parser/src/internal/infrastructure"
	"dc-aps-parser/src/pkg"
	"github.com/stretchr/testify/assert"
	"go.uber.org/goleak"
	"testing"
	"time"
)

func initAppWithMocks() (application.App, *output.NotificationAdapterMock, *output.TargetClientAdapterMock) {
	config := &infrastructure.Config{
		ParseInterval: time.Millisecond * 100,
	}

	targetClientAdapterMock := output.NewTargetClientAdapterMock()
	notificationAdapterMock := output.NewNotificationAdapterMock()
	resultService := application.NewResultService(targetClientAdapterMock)
	app := application.App{
		ResultService: resultService,
		ParserService: application.NewParserService(config, resultService, notificationAdapterMock),
		AdminService:  application.NewAdminService(),
	}
	return app, notificationAdapterMock, targetClientAdapterMock
}

func Test_Parser(t *testing.T) {
	t.Run("Parser creates", func(t *testing.T) {
		defer goleak.VerifyNone(t)

		app, notificationAdapterMock, targetClient := initAppWithMocks()
		targetClient.SetResults([]domain.ParseResult{
			{Items: make([]domain.ParseItem, 0)},
		})
		defer app.ParserService.StopAllParsersSync()

		_, err := app.ParserService.NewParser(1)

		if err != nil {
			t.Errorf("GetResult() error = %v", err)
			return
		}

		sentMessages := notificationAdapterMock.GetSentMessages()
		assert.Equal(t, 1, len(sentMessages))
		assert.Equal(t, "Парсер запущен", sentMessages[0].Text)
		//assert.Equal(t, "Найдено 0 объявлений. Ищу новые...", sentMessages[0].Text)
		assert.Equal(t, int64(1), sentMessages[0].ChatID)
	})

	t.Run("Parser doesn't stop if not exist", func(t *testing.T) {
		defer goleak.VerifyNone(t)
		app, _, _ := initAppWithMocks()
		defer app.ParserService.StopAllParsersSync()
		assert.Error(t, app.ParserService.StopParser(1))
	})

	t.Run("Many parsers create", func(t *testing.T) {
		defer goleak.VerifyNone(t)
		app, notificationAdapterMock, targetClient := initAppWithMocks()
		targetClient.SetResults([]domain.ParseResult{{Items: make([]domain.ParseItem, 0)}})
		defer app.ParserService.StopAllParsersSync()

		notificationAdapterMock.SetCalls(6)

		_, err := app.ParserService.NewParser(1)
		assert.NoError(t, err)
		_, err = app.ParserService.NewParser(2)
		assert.NoError(t, err)
		_, err = app.ParserService.NewParser(3)
		assert.NoError(t, err)

		notificationAdapterMock.WaitForCalls()

		sentMessages := notificationAdapterMock.GetSentMessages()
		assert.Equal(t, 6, len(sentMessages))

		check := func(chatId int64, messages []output.SentMessageMock) {
			messagesToChat := pkg.Filter(messages, func(msg output.SentMessageMock) bool {
				return msg.ChatID == chatId
			})
			assert.Equal(t, 2, len(messagesToChat))
			assert.Equal(t, "Парсер запущен", messagesToChat[0].Text)
			assert.Equal(t, "Найдено 0 объявлений. Ищу новые...", messagesToChat[1].Text)
		}

		check(1, sentMessages)
		check(2, sentMessages)
		check(3, sentMessages)
	})

	t.Run("Parser gets init aps count", func(t *testing.T) {
		defer goleak.VerifyNone(t)

		app, notificationAdapterMock, targetClient := initAppWithMocks()
		targetClient.SetResults([]domain.ParseResult{
			{Items: []domain.ParseItem{
				domain.NewParseItem(1, ""),
				domain.NewParseItem(2, ""),
				domain.NewParseItem(3, ""),
			}},
		})
		defer app.ParserService.StopAllParsersSync()

		notificationAdapterMock.SetCalls(2)
		_, err := app.ParserService.NewParser(1)

		if err != nil {
			t.Errorf("GetResult() error = %v", err)
			return
		}

		notificationAdapterMock.WaitForCalls()

		sentMessages := notificationAdapterMock.GetSentMessages()
		assert.Equal(t, "Парсер запущен", sentMessages[0].Text)
		assert.Equal(t, "Найдено 3 объявлений. Ищу новые...", sentMessages[1].Text)
	})
}
