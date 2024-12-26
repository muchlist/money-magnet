package fcm

import (
	"context"
	"fmt"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/messaging"
	"github.com/muchlist/moneymagnet/pkg/ds"
)

// FCMSender defines the interface for FCM client operations.
type FCMSender interface {
	SendMessage(ctx context.Context, payload Payload) ([]string /*invalid token*/, error)
}

type fcmClient struct {
	app    *firebase.App
	client *messaging.Client
}

func NewFcmClient(app *firebase.App) (FCMSender, error) {
	client, err := app.Messaging(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to get FCM client: %w", err)
	}
	return &fcmClient{app: app, client: client}, nil
}

// SendMessage sends a notification message to the specified receiver tokens.
func (c *fcmClient) SendMessage(ctx context.Context, payload Payload) ([]string /*invalid token*/, error) {
	// Validate receiver tokens
	validTokens, err := validateTokens(payload.ReceiverTokens)
	if err != nil {
		return nil, err
	}

	if len(validTokens) == 0 {
		return nil, nil
	}

	// Create message with the valid tokens
	message := createMulticastMessage(payload, validTokens)

	// Send message
	resp, err := c.client.SendEachForMulticast(ctx, message)
	if err != nil {
		return nil, fmt.Errorf("failed to send message: %w", err)
	}

	var invalidTokens []string
	if resp.FailureCount > 0 {
		for i, result := range resp.Responses {
			if result.Error != nil {
				invalidTokens = append(invalidTokens, validTokens[i])
			}
		}
	}

	return invalidTokens, nil
}

// validateTokens filters out invalid or empty tokens from the list.
func validateTokens(tokens []string) ([]string, error) {
	if len(tokens) == 0 {
		return nil, fmt.Errorf("receiver tokens cannot be empty or nil")
	}

	var validTokens []string
	for _, token := range tokens {
		if token != "" {
			validTokens = append(validTokens, token)
		}
	}

	return validTokens, nil
}

// createMulticastMessage creates a MulticastMessage struct with the given payload and tokens.
func createMulticastMessage(payload Payload, tokens []string) *messaging.MulticastMessage {
	uniqueToken := ds.NewStringSet()
	uniqueToken.AddAll(tokens)

	return &messaging.MulticastMessage{
		Notification: &messaging.Notification{
			Title: payload.Title,
			Body:  payload.Message,
		},
		Tokens: uniqueToken.RevealNotEmpty(),
	}
}
