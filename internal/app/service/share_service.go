package service

import (
	"context"
	"fmt"

	"go-tele-bot/internal/app/client/googledrive"
)

type ShareService struct {
	client *googledrive.GoogleDriveClient
}

// NewShareService creates a new ShareService
func NewShareService(ctx context.Context) (*ShareService, error) {
	client, err := googledrive.NewGoogleDriveClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create share service: %w", err)
	}
	return &ShareService{
		client: client,
	}, nil
}
