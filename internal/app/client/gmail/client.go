package gmail

import "context"

// TODO: Implement Gmail client
type GmailClient struct {
}

func NewGmailClient(ctx context.Context) (*GmailClient, error) {
	return &GmailClient{}, nil
}
