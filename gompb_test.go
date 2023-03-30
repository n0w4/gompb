package gompb_test

import (
	"testing"
	"time"

	"github.com/n0w4/gompb"
)

func TestNormalPath(t *testing.T) {
	pb := gompb.NewMemoryPubSub()

	handle := func(header map[string]interface{}, body []byte) {
		t.Run("expected event_type to be ChangeEmailRequested", func(t *testing.T) {
			want := "ChangeEmailRequested"
			got := header["event_type"]
			if want != got {
				t.Errorf("want %v but got %v", want, got)
			}
		})
	}

	poolSize := 1
	pb.Consume("test_topic", poolSize, handle)

	pb.Publish("test_topic", message())

	time.Sleep(2 * time.Second)
}

func message() gompb.Message {
	h := map[string]interface{}{
		"event_type": "ChangeEmailRequested",
		"command":    "ChangeEmail",
		"created_at": 1680067269,
		"message_id": "917d5819-86bc-47c9-98e2-4a7d192c15de",
	}

	body := `
	{
		"client_id": "17d5819-86bc-47c9-98e2-4a7d192c15de",
		"old_email": "client-1@oldprovider.com",
		"new_email": "client-1@newprovider.net"
	}`

	msg := gompb.Message{
		Header: h,
		Body:   []byte(body),
	}
	return msg
}
