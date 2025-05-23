package mocks

import (
	"context"

	"github.com/nazrawigedion123/wallet-backend/webhook/models"
)

type MockWebhookService struct {
	ProcessFunc func(ctx context.Context, payload models.IncomingWebhook) error
}

func (m *MockWebhookService) ProcessWebhook(ctx context.Context, payload models.IncomingWebhook) error {
	return m.ProcessFunc(ctx, payload)
}
