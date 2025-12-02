package googledrive

import (
	"context"
	"fmt"
	"os"

	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
)

// GoogleDriveClient is a client for the Google Drive API
type GoogleDriveClient struct {
	client *drive.Service
}

// NewGoogleDriveClient creates a new Google Drive client
func NewGoogleDriveClient(ctx context.Context) (*GoogleDriveClient, error) {
	var opts []option.ClientOption

	// Check for credentials JSON content in env var (useful for Cloud Run with Secret Manager)
	if credsJSON := os.Getenv("GOOGLE_CREDENTIALS_JSON"); credsJSON != "" {
		opts = append(opts, option.WithCredentialsJSON([]byte(credsJSON)))
	} else if credsFile := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS"); credsFile != "" {
		opts = append(opts, option.WithCredentialsFile(credsFile))
	} else {
		// Fallback to local credentials file for development
		opts = append(opts, option.WithCredentialsFile("credentials.json"))
	}

	client, err := drive.NewService(ctx, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create google drive client: %w", err)
	}
	return &GoogleDriveClient{
		client: client,
	}, nil
}

// ShareFolder shares the folder with the given email address using the Google Drive API
// Does not send email notification if using service account credentials
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
