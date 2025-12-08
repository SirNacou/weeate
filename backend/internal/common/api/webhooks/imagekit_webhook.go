package webhooks

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/imagekit-developer/imagekit-go/v2"
)

func RegisterImageKitWebhook(api huma.API, client *imagekit.Client) {
	handler := NewImageKitWebhookHandler(client)
	huma.Post(api, "/ik-webhook", handler.Handle, func(o *huma.Operation) {
		o.MaxBodyBytes = 65536 // 64KB
		o.Description = "Handle ImageKit webhooks for video transformations and uploads."
		o.Summary = "ImageKit Webhook Handler"
	})
}

type ImageKitWebhookHandler struct {
	client *imagekit.Client
}

func NewImageKitWebhookHandler(client *imagekit.Client) *ImageKitWebhookHandler {
	return &ImageKitWebhookHandler{
		client,
	}
}

func (h *ImageKitWebhookHandler) Handle(ctx context.Context, req *struct {
	Header  http.Header
	RawBody []byte
}) (*struct{}, error) {

	// Verify and unwrap webhook payload
	event, err := h.client.Webhooks.Unwrap(req.RawBody, req.Header)
	if err != nil {
		slog.ErrorContext(ctx, "Invalid webhook signature or malformed payload: %v", "error", err)
		return nil, huma.Error401Unauthorized("Invalid webhook signature or malformed payload")
	}

	slog.InfoContext(ctx, "Verified webhook event: %s", "type", event.Type)

	// Handle different event types with full type safety
	switch event.Type {
	case "video.transformation.accepted":
		videoEvent := event.AsVideoTransformationAcceptedEvent()
		slog.InfoContext(ctx, "Video transformation accepted: %s", "url", videoEvent.Data.Asset.URL)
		// Debugging: Track transformation requests
		// handleVideoTransformationAccepted(videoEvent)

	case "video.transformation.ready":
		videoEvent := event.AsVideoTransformationReadyEvent()
		slog.InfoContext(ctx, "Video transformation ready: %s", "url", videoEvent.Data.Transformation.Output.URL)
		// Update your database/CMS to show the transformed video
		// handleVideoTransformationReady(videoEvent)

	case "video.transformation.error":
		videoEvent := event.AsVideoTransformationErrorEvent()
		slog.InfoContext(ctx, "Video transformation error: %s", "reason", videoEvent.Data.Transformation.Error.Reason)
		// Log error and check your origin/URL endpoint settings
		// handleVideoTransformationError(videoEvent)

	case "upload.pre-transform.success":
		uploadEvent := event.AsUploadPreTransformSuccessEvent()
		slog.InfoContext(ctx, "Pre-transform success: %s", "file_id", uploadEvent.Data.FileID)
		// File uploaded and pre-transformation completed
		// handleUploadPreTransformSuccess(uploadEvent)

	case "upload.post-transform.success":
		postEvent := event.AsUploadPostTransformSuccessEvent()
		slog.InfoContext(ctx, "Post-transform success: %s", "name", postEvent.Data.Name)
		// Additional transformation completed
		// handleUploadPostTransformSuccess(postEvent)

	// Handle other event types as needed
	default:
		slog.InfoContext(ctx, "Unhandled event type: %s", "type", event.Type)
	}

	return &struct{}{}, nil
}
