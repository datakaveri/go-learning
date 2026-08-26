package application

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	dxsql "github.com/datakaveri/dx-common-go/platform/database/sql"
	platformerrors "github.com/datakaveri/dx-common-go/platform/errors"
	"github.com/datakaveri/dx-common-go/platform/security/identity"

	"github.com/datakaveri/dx-example-go/internal/domain"
)

type Repository interface {
	Create(context.Context, domain.Widget) error
	ByID(context.Context, string) (domain.Widget, error)
}

type Authorizer interface {
	Allowed(ctx context.Context, subject identity.Subject, resourceID, relation string) (bool, error)
}

type EventWriter interface {
	WidgetCreated(context.Context, domain.Widget, string) error
	Kick()
}

type Clock func() time.Time

type Service struct {
	repo   Repository
	authz  Authorizer
	events EventWriter
	tx     dxsql.Manager
	now    Clock
}

func New(repo Repository, authz Authorizer, events EventWriter, tx dxsql.Manager, now Clock) (*Service, error) {
	if repo == nil || authz == nil || events == nil || tx == nil || now == nil {
		return nil, fmt.Errorf("widget service: all dependencies are required")
	}
	return &Service{repo: repo, authz: authz, events: events, tx: tx, now: now}, nil
}

type CreateCommand struct {
	Subject       identity.Subject
	Name          string
	CorrelationID string
}

func (s *Service) Create(ctx context.Context, cmd CreateCommand) (domain.Widget, error) {
	allowed, err := s.authz.Allowed(ctx, cmd.Subject, "widgets", "editor")
	if err != nil {
		return domain.Widget{}, platformerrors.Wrap(err, platformerrors.CodeServiceUnavailable, "authorization unavailable")
	}
	if !allowed {
		return domain.Widget{}, platformerrors.Forbidden("access denied")
	}

	widget, err := domain.NewWidget(uuid.NewString(), cmd.Subject.ID, cmd.Subject.Org, cmd.Name, s.now())
	if err != nil {
		return domain.Widget{}, err
	}

	if err := s.tx.Do(ctx, func(txCtx context.Context) error {
		if err := s.repo.Create(txCtx, widget); err != nil {
			return err
		}
		return s.events.WidgetCreated(txCtx, widget, cmd.CorrelationID)
	}); err != nil {
		return domain.Widget{}, fmt.Errorf("create widget: %w", err)
	}
	s.events.Kick()
	return widget, nil
}

func (s *Service) Get(ctx context.Context, subject identity.Subject, id string) (domain.Widget, error) {
	allowed, err := s.authz.Allowed(ctx, subject, id, "viewer")
	if err != nil {
		return domain.Widget{}, platformerrors.Wrap(err, platformerrors.CodeServiceUnavailable, "authorization unavailable")
	}
	if !allowed {
		return domain.Widget{}, platformerrors.Forbidden("access denied")
	}
	widget, err := s.repo.ByID(ctx, id)
	if err != nil {
		return domain.Widget{}, fmt.Errorf("get widget %q: %w", id, err)
	}
	if widget.OrganisationID() != subject.Org {
		return domain.Widget{}, platformerrors.Forbidden("access denied")
	}
	return widget, nil
}
