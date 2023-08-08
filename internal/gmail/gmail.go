package gmail

import (
	"context"
	"fmt"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
	"net/http"
)

const user = "me"

type Gmail struct {
	svc    *gmail.Service
	client *http.Client
}

func New(ctx context.Context, client *http.Client) (*Gmail, error) {
	g := &Gmail{
		client: client,
	}
	svc, err := gmail.NewService(ctx, option.WithHTTPClient(g.client))
	if err != nil {
		return nil, fmt.Errorf("gmail: failed to init service: %w", err)
	}
	g.svc = svc
	return g, nil
}

func (g *Gmail) ListInbox() (*gmail.ListMessagesResponse, error) {
	return g.svc.Users.Messages.List(user).Q("in:inbox").Do()
}

func (g *Gmail) ListSpam() (*gmail.ListMessagesResponse, error) {
	return g.svc.Users.Messages.List(user).Q("in:spam").Do()
}

func (g *Gmail) GetMessage(id string) (*gmail.Message, error) {
	return g.svc.Users.Messages.Get(user, id).Do()
}

func (g *Gmail) DeleteMessage(id string) error {
	return g.svc.Users.Messages.Delete(user, id).Do()
}
