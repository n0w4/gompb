package gompb

import "sync"

type Header struct {
	EventType string
	Command   string
	CreatedAt int64
	MessageId string
}

type Message struct {
	Header Header
	Body   []byte
}

type MemoryPubSub struct {
	subscribers map[string][]chan Message
	mu          sync.RWMutex
}

func NewMemoryPubSub() *MemoryPubSub {
	return &MemoryPubSub{
		subscribers: make(map[string][]chan Message),
	}
}

func (ps *MemoryPubSub) Subscribe(topic string, pool int) <-chan Message {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	ch := make(chan Message, pool)
	ps.subscribers[topic] = append(ps.subscribers[topic], ch)
	return ch
}

func (ps *MemoryPubSub) Publish(topic string, msg Message) {
	ps.mu.RLock()
	defer ps.mu.RUnlock()

	for _, ch := range ps.subscribers[topic] {
		ch <- msg
	}
}

func (ps *MemoryPubSub) Unsubscribe(topic string, ch <-chan Message) {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	for i, sub := range ps.subscribers[topic] {
		if sub == ch {
			ps.subscribers[topic] = append(ps.subscribers[topic][:i], ps.subscribers[topic][i+1:]...)
			break
		}
	}
}

func (ps *MemoryPubSub) CloseTopic(topic string) error {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	for _, ch := range ps.subscribers[topic] {
		close(ch)
	}

	return nil
}

func (ps *MemoryPubSub) Close() error {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	for _, chs := range ps.subscribers {
		for _, ch := range chs {
			close(ch)
		}
	}

	return nil
}
