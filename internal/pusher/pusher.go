package pusher

import (
	"os"

	"github.com/pusher/pusher-http-go/v5"
)

func NewClient() pusher.Client {
	return pusher.Client{
		AppID:   os.Getenv("PUSHER_APP_ID"),
		Key:     os.Getenv("PUSHER_KEY"),
		Secret:  os.Getenv("PUSHER_SECRET"),
		Cluster: os.Getenv("PUSHER_CLUSTER"),
		Secure:  true,
	}
}
