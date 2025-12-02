package service

import (
	"context"
	"fmt"
	"strings"

	"go-tele-bot/internal/app/client/googledrive"
	"go-tele-bot/internal/data"
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

// ShareFolders shares the folders with the given email address
func (s *ShareService) ShareFolders(ctx context.Context, selected map[string]bool, email string) error {
	folderNames := make([]string, 0)
	for name := range selected {
		if selected[name] {
			folderNames = append(folderNames, name)
		}
	}
	folderIDs, err := data.GetFolderIDsByNames(folderNames)
	if err != nil {
		return fmt.Errorf("failed to get folder IDs by names: %w", err)
	}
	if len(folderIDs) == 0 {
		return fmt.Errorf("no folders to share")
	}

	// if email does not contain @, append @gmail.com
	if !strings.Contains(email, "@") {
		email = fmt.Sprintf("%s@gmail.com", email)
	}

	err = s.client.ShareFolders(ctx, folderIDs, email)
	if err != nil {
		return fmt.Errorf("failed to share folders: %w", err)
	}

	// TODO: send email to the user

	return nil
}
