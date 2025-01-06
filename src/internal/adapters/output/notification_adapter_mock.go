package output

import (
	"fmt"
	"sync"
)

type NotificationAdapterMock struct {
	sentMessages      []string
	wg                *sync.WaitGroup
	isWaitingForCalls bool
}

func NewNotificationAdapterMock() *NotificationAdapterMock {
	return &NotificationAdapterMock{
		sentMessages: make([]string, 0),
		wg:           &sync.WaitGroup{},
	}
}

func (n *NotificationAdapterMock) SendMessage(text string) error {
	n.sentMessages = append(n.sentMessages, text)
	fmt.Printf("Sent mock message: %s\n", text)
	if n.isWaitingForCalls {
		n.wg.Done()
	}
	return nil
}

func (n *NotificationAdapterMock) GetSentMessages() []string {
	return n.sentMessages
}

func (n *NotificationAdapterMock) SetCalls(i int) {
	n.isWaitingForCalls = true
	n.wg.Add(i)
}

func (n *NotificationAdapterMock) WaitForCalls() {
	n.wg.Wait()
}
