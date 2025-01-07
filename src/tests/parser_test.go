package tests

import (
	"dc-aps-parser/src/internal/adapters/output"
	"dc-aps-parser/src/internal/core/application"
	"dc-aps-parser/src/internal/infrastructure"
	"dc-aps-parser/src/pkg"
	"fmt"
	"github.com/stretchr/testify/assert"
	"go.uber.org/goleak"
	"testing"
	"time"
)

func initAppWithMocks() (application.App, *output.NotificationAdapterMock, *output.TargetClientAdapterMock, *infrastructure.Config) {
	config := &infrastructure.Config{
		ParseInterval:                 time.Millisecond * 100,
		TgAdminChatId:                 int64(10),
		TgParserLaunchMessage:         "Parser launched",
		TgUserStartMessage:            "Hello!",
		TgParserAlreadyStoppedMessage: "Parser already stopped",
		TgErrorStoppingParserMessage:  "Error stopping parser",
		TgParserStoppedMessage:        "Parser stopped",
		TgAdminHelpMessage:            "help admin",
		TgUserHelpMessage:             "help user",
		TgStoppedParserStatusMessage:  "parser is stopped",
		TgUnknownCommandMessage:       "Unknown command",
		TgActiveParserStatus:          "Parser is active",
		TgInitialApsCountFormat:       "init %d aps",
	}

	targetClientAdapterMock := output.NewTargetClientAdapterMock()
	notificationAdapterMock := output.NewNotificationAdapterMock()
	resultService := application.NewResultService(targetClientAdapterMock)
	parserNotificationService := application.NewParserNotificationService(config, notificationAdapterMock)
	app := application.App{
		ResultService: resultService,
		ParserService: application.NewParserService(config, resultService, parserNotificationService, nil),
		AdminService:  application.NewAdminService(config),
	}
	return app, notificationAdapterMock, targetClientAdapterMock, config
}

func Test_ParserLaunch(t *testing.T) {
	t.Run("Parser creates", func(t *testing.T) {
		defer goleak.VerifyNone(t)

		app, notificationAdapterMock, _, config := initAppWithMocks()
		defer app.ParserService.StopAllParsersSync()

		_, err := app.ParserService.NewParser(1, "", false)

		if err != nil {
			t.Errorf("GetResult() error = %v", err)
			return
		}

		sentMessages := notificationAdapterMock.GetSentMessages()
		assert.Equal(t, 1, len(sentMessages))
		assert.Equal(t, config.TgParserLaunchMessage, sentMessages[0].Text)
		assert.Equal(t, int64(1), sentMessages[0].ChatID)
	})

	t.Run("Parser doesn't stop if not exist", func(t *testing.T) {
		defer goleak.VerifyNone(t)
		app, _, _, _ := initAppWithMocks()
		defer app.ParserService.StopAllParsersSync()
		assert.Error(t, app.ParserService.StopParser(1))
	})

	t.Run("Many parsers create", func(t *testing.T) {
		defer goleak.VerifyNone(t)
		app, notificationAdapterMock, _, config := initAppWithMocks()
		defer app.ParserService.StopAllParsersSync()

		notificationAdapterMock.SetCalls(6)

		_, err := app.ParserService.NewParser(1, "", false)
		assert.NoError(t, err)
		_, err = app.ParserService.NewParser(2, "", false)
		assert.NoError(t, err)
		_, err = app.ParserService.NewParser(3, "", false)
		assert.NoError(t, err)

		notificationAdapterMock.WaitForCalls()

		sentMessages := notificationAdapterMock.GetSentMessages()
		assert.Equal(t, 6, len(sentMessages))

		check := func(chatId int64, messages []output.SentMessageMock) {
			messagesToChat := pkg.Filter(messages, func(msg output.SentMessageMock) bool {
				return msg.ChatID == chatId
			})
			assert.Equal(t, 2, len(messagesToChat))
			assert.Equal(t, config.TgParserLaunchMessage, messagesToChat[0].Text)
			assert.Equal(t, fmt.Sprintf(config.TgInitialApsCountFormat, 0), messagesToChat[1].Text)
		}

		check(1, sentMessages)
		check(2, sentMessages)
		check(3, sentMessages)
	})

	t.Run("Parser gets init aps count", func(t *testing.T) {
		defer goleak.VerifyNone(t)

		app, notificationAdapterMock, targetClient, config := initAppWithMocks()
		targetClient.SetResults([]int{3})
		defer app.ParserService.StopAllParsersSync()

		notificationAdapterMock.SetCalls(2)
		_, err := app.ParserService.NewParser(1, "", false)

		if err != nil {
			t.Errorf("GetResult() error = %v", err)
			return
		}

		notificationAdapterMock.WaitForCalls()

		sentMessages := notificationAdapterMock.GetSentMessages()
		assert.Equal(t, config.TgParserLaunchMessage, sentMessages[0].Text)
		assert.Equal(t, fmt.Sprintf(config.TgInitialApsCountFormat, 3), sentMessages[1].Text)
	})
}

func TestParserWorks(t *testing.T) {
	t.Run("Parser detects new ap", func(t *testing.T) {
		defer goleak.VerifyNone(t)
		app, notificationAdapter, targetClient, _ := initAppWithMocks()
		defer app.ParserService.StopAllParsersSync()

		targetClient.SetResults([]int{10, 10, 10, 11})
		notificationAdapter.SetCalls(3)
		_, err := app.ParserService.NewParser(1, "", false)
		assert.NoError(t, err)
		notificationAdapter.WaitForCalls()

		sentMessages := notificationAdapter.GetSentMessages()
		assert.Equal(t, "Новое объявление link_11", sentMessages[2].Text)
	})

	t.Run("Parser detects many new aps", func(t *testing.T) {
		defer goleak.VerifyNone(t)
		app, notificationAdapter, targetClient, _ := initAppWithMocks()
		defer app.ParserService.StopAllParsersSync()

		targetClient.SetResults([]int{10, 10, 11, 11, 12})
		notificationAdapter.SetCalls(4)
		_, err := app.ParserService.NewParser(1, "", false)
		assert.NoError(t, err)
		notificationAdapter.WaitForCalls()

		sentMessages := notificationAdapter.GetSentMessages()
		assert.Equal(t, "Новое объявление link_11", sentMessages[2].Text)
		assert.Equal(t, "Новое объявление link_12", sentMessages[3].Text)
	})

	t.Run("Parser detects many new aps at one time", func(t *testing.T) {
		defer goleak.VerifyNone(t)
		app, notificationAdapter, targetClient, _ := initAppWithMocks()
		defer app.ParserService.StopAllParsersSync()

		targetClient.SetResults([]int{10, 14})
		notificationAdapter.SetCalls(6)
		_, err := app.ParserService.NewParser(1, "", false)
		assert.NoError(t, err)
		notificationAdapter.WaitForCalls()

		sentMessages := notificationAdapter.GetSentMessages()
		assert.Equal(t, 6, len(sentMessages))
	})

	t.Run("Parser ignores hiding and showing some ap", func(t *testing.T) {
		defer goleak.VerifyNone(t)
		app, notificationAdapter, targetClient, _ := initAppWithMocks()
		defer app.ParserService.StopAllParsersSync()

		targetClient.SetResults([]int{10, 9, 10, 11})
		notificationAdapter.SetCalls(3)
		_, err := app.ParserService.NewParser(1, "", false)
		assert.NoError(t, err)
		notificationAdapter.WaitForCalls()

		sentMessages := notificationAdapter.GetSentMessages()
		assert.Equal(t, "Новое объявление link_11", sentMessages[2].Text)
	})
}
