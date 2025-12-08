package webhooks

import (
	"context"
	"fmt"
	"net/http"
	"os"

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
		fmt.Fprintf(os.Stderr, "Invalid webhook signature or malformed payload: %v\n", err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	fmt.Printf("Verified webhook event: %s\n", event.Type)

	// Handle different event types with full type safety
	switch event.Type {
	case "video.transformation.accepted":
		videoEvent := event.AsVideoTransformationAcceptedEvent()
		fmt.Printf("Video transformation accepted: %s\n", videoEvent.Data.Asset.URL)
		// Debugging: Track transformation requests
		// handleVideoTransformationAccepted(videoEvent)

	case "video.transformation.ready":
		videoEvent := event.AsVideoTransformationReadyEvent()
		fmt.Printf("Video transformation ready: %s\n", videoEvent.Data.Transformation.Output.URL)
		// Update your database/CMS to show the transformed video
		// handleVideoTransformationReady(videoEvent)

	case "video.transformation.error":
		videoEvent := event.AsVideoTransformationErrorEvent()
		fmt.Printf("Video transformation error: %s\n", videoEvent.Data.Transformation.Error.Reason)
		// Log error and check your origin/URL endpoint settings
		// handleVideoTransformationError(videoEvent)

	case "upload.pre-transform.success":
		uploadEvent := event.AsUploadPreTransformSuccessEvent()
		fmt.Printf("Pre-transform success: %s\n", uploadEvent.Data.FileID)
		// File uploaded and pre-transformation completed
		// handleUploadPreTransformSuccess(uploadEvent)

	case "upload.post-transform.success":
		postEvent := event.AsUploadPostTransformSuccessEvent()
		fmt.Printf("Post-transform success: %s\n", postEvent.Data.Name)
		// Additional transformation completed
		// handleUploadPostTransformSuccess(postEvent)

	// Handle other event types as needed
	default:
		fmt.Printf("Unhandled event type: %s\n", event.Type)
	}

	w.WriteHeader(http.StatusOK)
	return nil, nil
}
