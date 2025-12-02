package googledrive

import (
	"context"
	"fmt"

	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
)

// GoogleDriveClient is a client for the Google Drive API
type GoogleDriveClient struct {
	client *drive.Service
}

// NewGoogleDriveClient creates a new Google Drive client
func NewGoogleDriveClient(ctx context.Context) (*GoogleDriveClient, error) {
	client, err := drive.NewService(ctx, option.WithCredentialsFile("credentials.json"))
	if err != nil {
		return nil, fmt.Errorf("failed to create google drive client: %w", err)
	}
	return &GoogleDriveClient{
		client: client,
	}, nil
}

// ShareFolder shares the folder with the given email address using the Google Drive API
func (c *GoogleDriveClient) ShareFolder(ctx context.Context, folderID string, email string) error {
	_, err := c.client.Permissions.Create(folderID, &drive.Permission{
		Role:         "reader",
		Type:         "user",
		EmailAddress: email,
	}).Do()
	if err != nil {
		return fmt.Errorf("failed to share folder with email %s: %w", email, err)
	}
	return nil
}

// ShareFolders shares the folders with the given email address using the Google Drive API
func (c *GoogleDriveClient) ShareFolders(ctx context.Context, folderIDs []string, email string) error {
	for _, folderID := range folderIDs {
		err := c.ShareFolder(ctx, folderID, email)
		if err != nil {
			return fmt.Errorf("failed to share folder %s with email %s: %w", folderID, email, err)
		}
	}
	return nil
}
