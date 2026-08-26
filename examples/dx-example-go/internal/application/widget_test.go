package application

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	dxsql "github.com/datakaveri/dx-common-go/platform/database/sql"
	platformerrors "github.com/datakaveri/dx-common-go/platform/errors"
	"github.com/datakaveri/dx-common-go/platform/security/identity"

	"github.com/datakaveri/dx-example-go/internal/domain"
)

type fakeRepository struct {
	createCalls int
	widget      domain.Widget
}

func (f *fakeRepository) Create(_ context.Context, widget domain.Widget) error {
	f.createCalls++
	f.widget = widget
	return nil
}

func (f *fakeRepository) ByID(context.Context, string) (domain.Widget, error) { return f.widget, nil }

type fakeAuthorizer struct{ allowed bool }

func (f fakeAuthorizer) Allowed(context.Context, identity.Subject, string, string) (bool, error) {
	return f.allowed, nil
}

type fakeEvents struct{ writes int }

func (f *fakeEvents) WidgetCreated(context.Context, domain.Widget, string) error {
	f.writes++
	return nil
}
func (*fakeEvents) Kick() {}

type passthroughTx struct{}

func (passthroughTx) Do(ctx context.Context, fn func(context.Context) error, _ ...dxsql.TxOption) error {
	return fn(ctx)
}
func (passthroughTx) DoRetry(ctx context.Context, fn func(context.Context) error, opts ...dxsql.TxOption) error {
	return passthroughTx{}.Do(ctx, fn, opts...)
}
func (passthroughTx) Lock(ctx context.Context, _ string, fn func(context.Context) error) error {
	return fn(ctx)
}

func TestCreateDeniedDoesNotWrite(t *testing.T) {
	t.Parallel()

	repo := &fakeRepository{}
	eventWriter := &fakeEvents{}
	service, err := New(repo, fakeAuthorizer{allowed: false}, eventWriter, passthroughTx{}, func() time.Time {
		return time.Date(2026, time.August, 10, 0, 0, 0, 0, time.UTC)
	})
	require.NoError(t, err)

	_, err = service.Create(context.Background(), CreateCommand{
		Subject: identity.Subject{ID: "user-1", Org: "org-1"}, Name: "Weather",
	})

	require.ErrorIs(t, err, platformerrors.ErrForbidden)
	require.Zero(t, repo.createCalls)
	require.Zero(t, eventWriter.writes)
}

func TestCreateWritesStateAndEvent(t *testing.T) {
	t.Parallel()

	repo := &fakeRepository{}
	eventWriter := &fakeEvents{}
	service, err := New(repo, fakeAuthorizer{allowed: true}, eventWriter, passthroughTx{}, func() time.Time {
		return time.Date(2026, time.August, 10, 0, 0, 0, 0, time.UTC)
	})
	require.NoError(t, err)

	widget, err := service.Create(context.Background(), CreateCommand{
		Subject: identity.Subject{ID: "user-1", Org: "org-1"}, Name: " Weather ",
		CorrelationID: "request-1",
	})

	require.NoError(t, err)
	require.Equal(t, "Weather", widget.Name())
	require.Equal(t, 1, repo.createCalls)
	require.Equal(t, 1, eventWriter.writes)
}
