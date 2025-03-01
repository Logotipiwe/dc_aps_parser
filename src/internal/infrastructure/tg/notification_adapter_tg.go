package tg

type NotificationAdapterTg struct {
	bot *BotAPI
}

func NewNotificationAdapterTg(
	botAPI *BotAPI,
) *NotificationAdapterTg {
	return &NotificationAdapterTg{
		botAPI,
	}
}

func (n *NotificationAdapterTg) SendMessage(chatID int64, text string) error {
	err := n.bot.SendMessageInTg(chatID, text)
	return err
}

func (n *NotificationAdapterTg) SendMessageWithImages(chatID int64, text string, images []string) error {
	return n.bot.SendMessageInTgWithImages(chatID, text, images)
}
