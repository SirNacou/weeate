package imagekit

import (
	"github.com/imagekit-developer/imagekit-go/v2"
	"github.com/imagekit-developer/imagekit-go/v2/option"
)

func NewImageKitClient(privateKey string, webhookSecret string) *imagekit.Client {
	c := imagekit.NewClient(
		option.WithPrivateKey(privateKey),
		option.WithWebhookSecret(webhookSecret),
	)

	return &c
}
