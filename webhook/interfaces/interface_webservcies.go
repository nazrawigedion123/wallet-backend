
package interfaces

import (
	"context"
	"github.com/nazrawigedion123/wallet-backend/webhook/models"
)




type WebhookService interface {
	ProcessWebhook(ctx context.Context, payload models.IncomingWebhook) error
}