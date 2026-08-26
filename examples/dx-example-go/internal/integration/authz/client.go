package authz

import (
	"context"

	"github.com/datakaveri/dx-common-go/auth/fga"
	"github.com/datakaveri/dx-common-go/platform/security/identity"
)

type Client struct{ fga *fga.Client }

func New(client *fga.Client) *Client { return &Client{fga: client} }

func (c *Client) Allowed(ctx context.Context, subject identity.Subject, resourceID, relation string) (bool, error) {
	decision, err := c.fga.Check(ctx, fga.CheckRequest{
		SubjectType:  fga.SubjectTypeUser,
		SubjectID:    subject.ID,
		ResourceType: "resource",
		ResourceID:   resourceID,
		Relation:     relation,
	})
	if err != nil {
		return false, err
	}
	return decision.Allowed, nil
}
