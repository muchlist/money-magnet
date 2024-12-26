package fcm

import (
	"context"
	"fmt"

	firebase "firebase.google.com/go/v4"
	"google.golang.org/api/option"
)

// Config holds the configuration for FCM initialization.
type Config struct {
	CredLocation string
}

// InitFirebase initializes the Firebase app for sending notifications.
func InitFirebase(cfg Config) (*firebase.App, error) {
	opt := option.WithCredentialsFile(cfg.CredLocation)

	// Initialize Firebase App
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize Firebase app: %w", err)
	}

	return app, nil
}
