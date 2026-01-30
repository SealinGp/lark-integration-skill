package larkclient

import (
	"context"

	lark "github.com/larksuite/oapi-sdk-go/v3"
	larkcore "github.com/larksuite/oapi-sdk-go/v3/core"
)

type ClientWrapper struct {
	Client *lark.Client
}

func NewClient(appID, appSecret string) *ClientWrapper {
	// Initialize Lark Client with recommended settings
	client := lark.NewClient(appID, appSecret,
		lark.WithLogReqAtDebug(true),
		lark.WithLogLevel(larkcore.LogLevelInfo),
	)
	return &ClientWrapper{Client: client}
}

// Helper to get Tenant Access Token context (if needed explicitly, 
// though SDK usually handles token management automatically for default requests)
func (c *ClientWrapper) GetContext() context.Context {
	return context.Background()
}
