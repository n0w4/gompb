package gompb

import "sync"

type Message struct {
	Header map[string]interface{}
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

func (ps *MemoryPubSub) subscribe(topic string, pool int) <-chan Message {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	ch := make(chan Message, pool)
	ps.subscribers[topic] = append(ps.subscribers[topic], ch)
	return ch
}

func (ps *MemoryPubSub) Publish(topic string, msg Message) bool {
	ps.mu.RLock()
	defer ps.mu.RUnlock()
	if _, ok := ps.subscribers[topic]; !ok {
		return false
	}
	for _, ch := range ps.subscribers[topic] {
		ch <- msg
	}
	return true
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

func (ps *MemoryPubSub) Subscribers(topic string) int {
	ps.mu.RLock()
	defer ps.mu.RUnlock()

	return len(ps.subscribers[topic])
}

func (ps *MemoryPubSub) Topics() []string {
	ps.mu.RLock()
	defer ps.mu.RUnlock()

	var topics []string
	for topic := range ps.subscribers {
		topics = append(topics, topic)
	}

	return topics
}

func (ps *MemoryPubSub) TopicExists(topic string) bool {
	ps.mu.RLock()
	defer ps.mu.RUnlock()

	_, ok := ps.subscribers[topic]
	return ok
}

func (ps *MemoryPubSub) TopicSize(topic string) int {
	ps.mu.RLock()
	defer ps.mu.RUnlock()

	return len(ps.subscribers[topic])
}

func (ps *MemoryPubSub) Consume(topic string, pool int, handler func(h map[string]interface{}, b []byte)) {
	ch := ps.subscribe(topic, pool)
	for i := 0; i < pool; i++ {
		go func() {
			for msg := range ch {
				handler(msg.Header, msg.Body)
			}
		}()
	}
}
