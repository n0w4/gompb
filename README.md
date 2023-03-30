# GOMPB

is a very simple in memory pub-sub 

## Example building a message

```go
// HEADER

h := map[string]interface{}{
    "event_type": "ChangeEmailRequested",
    "command":    "ChangeEmail",
    "created_at": 1680067269,
    "message_id": "917d5819-86bc-47c9-98e2-4a7d192c15de",
}
```
```go
// BODY

body := `{
    "client_id": "2513683e-90c8-4858-b4b3-3f8d3e129156",
    "old_email": "some.client@oldprovider.com",
    "new_email": "some_client@newprovider.net"
}`
```
```go
// MESSAGE

msg := gompb.Message{
    Header: h,
    Body: []byte(body)
}
```

## Example usage

- poolSize here define how many goroutines consume create
- when consume you subscribe on a topic 
- Publish return `false` when no have subscribers on topic, you need consume a topic before publish

```go
// SOME HANDLER

type MessageHandler struct{}

func (sh *MessageHandler) Handle(header map[string]interface{}, body []byte) {
	// do some thing 
}
```

```go
// ALL TOGETHER

handler := MessageHandler{}

pb := gompb.NewMemoryPubSub()

poolSize := 100
pb.Consume("test_topic", poolSize, handler)

pb.Publish("test_topic", msg1)
pb.Publish("test_topic", msg2)
pb.Publish("test_topic", msg3)
pb.Publish("test_topic", msg4)
```