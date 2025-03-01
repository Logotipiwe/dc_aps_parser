package mock

import (
	"fmt"
	"sync"
)

type NotificationAdapterMock struct {
	sentMessages      []SentMessageMock
	wg                *sync.WaitGroup
	isWaitingForCalls bool
}

type SentMessageMock struct {
	ChatID int64
	Text   string
	Images []string
}

func NewNotificationAdapterMock() *NotificationAdapterMock {
	return &NotificationAdapterMock{
		sentMessages: make([]SentMessageMock, 0),
		wg:           &sync.WaitGroup{},
	}
}

func (n *NotificationAdapterMock) SendMessage(chatID int64, text string) error {
	n.sentMessages = append(n.sentMessages, SentMessageMock{chatID, text, []string{}})
	fmt.Printf("Sent mock message: %s\n", text)
	if n.isWaitingForCalls {
		n.wg.Done()
	}
	return nil
}

func (n *NotificationAdapterMock) SendMessageWithImages(chatID int64, text string, images []string) error {
	n.sentMessages = append(n.sentMessages, SentMessageMock{chatID, text, images})
	fmt.Printf("Sent mock message with images: %s %v\n", text, images)
	if n.isWaitingForCalls {
		n.wg.Done()
	}
	return nil
}

func (n *NotificationAdapterMock) GetSentMessages() []SentMessageMock {
	return n.sentMessages
}

func (n *NotificationAdapterMock) SetCalls(i int) {
	n.isWaitingForCalls = true
	n.wg.Add(i)
}

func (n *NotificationAdapterMock) WaitForCalls() {
	n.wg.Wait()
}
