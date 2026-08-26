package http

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	dxsql "github.com/datakaveri/dx-common-go/platform/database/sql"
	httpx "github.com/datakaveri/dx-common-go/platform/http"
	"github.com/datakaveri/dx-common-go/platform/security/identity"

	"github.com/datakaveri/dx-example-go/internal/application"
	"github.com/datakaveri/dx-example-go/internal/domain"
)

type repositoryFake struct{ widget domain.Widget }

func (f *repositoryFake) Create(_ context.Context, widget domain.Widget) error {
	f.widget = widget
	return nil
}
func (f *repositoryFake) ByID(context.Context, string) (domain.Widget, error) { return f.widget, nil }

type authzFake struct{}

func (authzFake) Allowed(context.Context, identity.Subject, string, string) (bool, error) {
	return true, nil
}

type eventsFake struct{}

func (eventsFake) WidgetCreated(context.Context, domain.Widget, string) error { return nil }
func (eventsFake) Kick()                                                      {}

type txFake struct{}

func (txFake) Do(ctx context.Context, fn func(context.Context) error, _ ...dxsql.TxOption) error {
	return fn(ctx)
}
func (txFake) DoRetry(ctx context.Context, fn func(context.Context) error, opts ...dxsql.TxOption) error {
	return txFake{}.Do(ctx, fn, opts...)
}
func (txFake) Lock(ctx context.Context, _ string, fn func(context.Context) error) error {
	return fn(ctx)
}

func TestCreate(t *testing.T) {
	t.Parallel()

	app, err := application.New(&repositoryFake{}, authzFake{}, eventsFake{}, txFake{}, func() time.Time {
		return time.Date(2026, time.August, 10, 0, 0, 0, 0, time.UTC)
	})
	require.NoError(t, err)
	handler := NewHandler(app)

	got, err := handler.Create(context.Background(), CreateRequest{
		Actor: httpActor("user-1", "org-1"), Name: "Rain gauge",
	})
	require.NoError(t, err)
	require.Equal(t, "Rain gauge", got.Value.Name)
}

func httpActor(userID, orgID string) httpx.Actor {
	return httpx.Actor{Subject: identity.Subject{ID: userID, Org: orgID}}
}
