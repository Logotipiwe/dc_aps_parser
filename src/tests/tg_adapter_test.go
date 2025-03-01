package tests

import (
	"dc-aps-parser/src/internal/adapters"
	"dc-aps-parser/src/internal/core/application"
	. "dc-aps-parser/src/internal/core/domain"
	"dc-aps-parser/src/internal/infrastructure/mock"
	"dc-aps-parser/src/pkg"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/stretchr/testify/assert"
	"go.uber.org/goleak"
	"testing"
)

func createTgAdapter(app application.App) *adapters.ParserAdapterTg {
	return adapters.NewParserAdapterTg(nil, app.ParserService, app.AdminService, app.ParserNotificationService)
}

func newUpdate(chatID int64, text string, username string) tgbotapi.Update {
	return tgbotapi.Update{
		Message: &tgbotapi.Message{
			Text: text,
			Chat: &tgbotapi.Chat{
				ID: chatID,
			},
			From: &tgbotapi.User{
				UserName: username,
			},
		},
	}
}

func TestCommands(t *testing.T) {
	t.Run("Start answers", func(t *testing.T) {
		app, adapterMocks, config := initAppWithMocks()
		adapterTg := createTgAdapter(app)

		err := adapterTg.HandleTgUpdate(newUpdate(1, "/start", ""))
		assert.Nil(t, err)
		sentMessages := adapterMocks.notification.GetSentMessages()
		assert.Equal(t, 1, len(sentMessages))
		assert.Equal(t, config.TgUserStartMessage, sentMessages[0].Text)
	})

	t.Run("Stop command", func(t *testing.T) {
		t.Run("Already stopped", func(t *testing.T) {
			app, adapterMocks, config := initAppWithMocks()
			adapterTg := createTgAdapter(app)

			err := adapterTg.HandleTgUpdate(newUpdate(1, "/stop", ""))
			assert.Nil(t, err)
			sentMessages := adapterMocks.notification.GetSentMessages()
			assert.Equal(t, 1, len(sentMessages))
			assert.Equal(t, config.TgParserAlreadyStoppedMessage, sentMessages[0].Text)
		})
		t.Run("Stopped", func(t *testing.T) {
			defer goleak.VerifyNone(t)
			app, adapterMocks, config := initAppWithMocks()
			adapterTg := createTgAdapter(app)
			defer app.ParserService.StopAllParsersSync()

			_, err := app.ParserService.LaunchParser(ParserParams{ChatID: 1, ParseLink: "some", UserName: "username"})
			assert.Nil(t, err)

			err = adapterTg.HandleTgUpdate(newUpdate(1, "/stop", ""))
			assert.Nil(t, err)
			sentMessages := adapterMocks.notification.GetSentMessages()
			assert.True(t, pkg.Some(sentMessages, func(msg mock.SentMessageMock) bool {
				return msg.Text == config.TgParserStoppedMessage
			}))
		})
	})
	t.Run("Help user", func(t *testing.T) {
		app, adapterMocks, config := initAppWithMocks()
		adapterTg := createTgAdapter(app)

		err := adapterTg.HandleTgUpdate(newUpdate(1, "/help", ""))
		assert.Nil(t, err)
		sentMessages := adapterMocks.notification.GetSentMessages()
		assert.Equal(t, config.TgUserHelpMessage, sentMessages[0].Text)
	})
	t.Run("Help admin", func(t *testing.T) {
		app, adapterMocks, config := initAppWithMocks()
		adapterTg := createTgAdapter(app)

		err := adapterTg.HandleTgUpdate(newUpdate(config.TgAdminChatId, "/help", ""))
		assert.Nil(t, err)
		sentMessages := adapterMocks.notification.GetSentMessages()
		assert.Equal(t, config.TgAdminHelpMessage, sentMessages[0].Text)
	})
	t.Run("Info for admin", func(t *testing.T) {
		// TODO
	})
	t.Run("Active parser status", func(t *testing.T) {
		defer goleak.VerifyNone(t)
		app, adapterMocks, config := initAppWithMocks()
		adapterTg := createTgAdapter(app)
		defer app.ParserService.StopAllParsersSync()

		_, err := app.ParserService.LaunchParser(ParserParams{ChatID: 1, ParseLink: "some", UserName: "username"})
		assert.Nil(t, err)

		err = adapterTg.HandleTgUpdate(newUpdate(1, "/status", ""))
		assert.Nil(t, err)
		sentMessages := adapterMocks.notification.GetSentMessages()
		assert.True(t, pkg.Some(sentMessages, func(msg mock.SentMessageMock) bool {
			return msg.Text == config.TgActiveParserStatus
		}))
	})
	t.Run("Non active parser status", func(t *testing.T) {
		app, adapterMocks, config := initAppWithMocks()
		adapterTg := createTgAdapter(app)

		err := adapterTg.HandleTgUpdate(newUpdate(1, "/status", ""))
		assert.Nil(t, err)
		sentMessages := adapterMocks.notification.GetSentMessages()
		assert.Equal(t, config.TgStoppedParserStatusMessage, sentMessages[0].Text)
	})
}

func TestEnablingParser(t *testing.T) {
	t.Run("Unknown if wrong link", func(t *testing.T) {
		defer goleak.VerifyNone(t)
		app, adapterMocks, config := initAppWithMocks()
		adapterTg := createTgAdapter(app)
		adapterMocks.targetClient.SetResults([]int{4})
		defer app.ParserService.StopAllParsersSync()

		err := adapterTg.HandleTgUpdate(newUpdate(1, "https://other.site", ""))
		assert.Nil(t, err)

		sentMessages := adapterMocks.notification.GetSentMessages()
		assert.Equal(t, 1, len(sentMessages))
		assert.Equal(t, config.TgUnknownCommandMessage, sentMessages[0].Text)
	})

	t.Run("Starts if correct link", func(t *testing.T) {
		defer goleak.VerifyNone(t)
		app, adapterMocks, config := initAppWithMocks()
		adapterTg := createTgAdapter(app)
		adapterMocks.targetClient.SetResults([]int{4})
		adapterMocks.notification.SetCalls(2)
		defer app.ParserService.StopAllParsersSync()

		err := adapterTg.HandleTgUpdate(newUpdate(1, "https://www.avito.ru/js/1/map/items?kek", ""))
		assert.Nil(t, err)

		adapterMocks.notification.WaitForCalls()

		sentMessages := adapterMocks.notification.GetSentMessages()
		assert.Equal(t, config.TgParserLaunchMessage, sentMessages[0].Text)
		assert.Equal(t, fmt.Sprintf(config.TgInitialApsCountFormat, 4), sentMessages[1].Text)
	})
}
