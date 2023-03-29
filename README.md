# GOMPB

is a very simple in memory pub-sub 

## Example building a message

```go
// HEADER

type Header struct {
	EventType string
	Command   string
	CreatedAt int64
	MessageId string
}
h := Header{
    EventType: "ChangeEmailRequested",
    Command: "ChangeEmail",
    CreatedAt: 1680067269,
    MessageId: "917d5819-86bc-47c9-98e2-4a7d192c15de",
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