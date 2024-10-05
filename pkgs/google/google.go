package google

import (
	"context"

	"google.golang.org/api/chat/v1"
)

func init_client() {
	ctx := context.Background()
	chatService, err := chat.NewService(ctx)
}
