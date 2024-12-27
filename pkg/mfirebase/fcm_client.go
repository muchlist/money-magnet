package mfirebase

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
func (c *fcmClient) SendMessage(ctx context.Context, payload Payload) ([]string, error) {
	validTokens, err := validateTokens(payload.ReceiverTokens)
	if err != nil {
		return nil, err
	}

	if len(validTokens) == 0 {
		return nil, nil
	}

	var invalidTokens []string

	// Split tokens into chunks if exceeding max limit
	tokenChunks := chunkTokens(validTokens, 500)

	for _, tokens := range tokenChunks {
		message := createMulticastMessage(payload, tokens)
		resp, err := c.client.SendEachForMulticast(ctx, message)
		if err != nil {
			return nil, fmt.Errorf("failed to send message: %w", err)
		}

		if resp.FailureCount > 0 {
			invalid := extractInvalidTokens(tokens, resp.Responses)
			invalidTokens = append(invalidTokens, invalid...)
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

func chunkTokens(tokens []string, chunkSize int) [][]string {
	var chunks [][]string
	for i := 0; i < len(tokens); i += chunkSize {
		end := i + chunkSize
		if end > len(tokens) {
			end = len(tokens)
		}
		chunks = append(chunks, tokens[i:end])
	}
	return chunks
}

func extractInvalidTokens(tokens []string, responses []*messaging.SendResponse) []string {
	var invalidTokens []string
	for i, result := range responses {
		if result.Error != nil {
			invalidTokens = append(invalidTokens, tokens[i])
		}
	}
	return invalidTokens
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
