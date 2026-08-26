package postgres

import (
	"context"
	"fmt"
	"time"

	dxsql "github.com/datakaveri/dx-common-go/platform/database/sql"

	"github.com/datakaveri/dx-example-go/internal/domain"
)

type widgetRow struct {
	ID             string    `db:"id"`
	OwnerID        string    `db:"owner_id"`
	OrganisationID string    `db:"organisation_id"`
	Name           string    `db:"name"`
	CreatedAt      time.Time `db:"created_at"`
}

type WidgetRepository struct{ rows *dxsql.Repo[widgetRow] }

func NewWidgetRepository(db dxsql.DB) *WidgetRepository {
	return &WidgetRepository{rows: dxsql.NewRepo[widgetRow](db, dxsql.WithTable[widgetRow]("widgets"))}
}

func (r *WidgetRepository) Create(ctx context.Context, widget domain.Widget) error {
	_, err := r.rows.Insert(ctx, map[dxsql.Column]any{
		"id": widget.ID(), "owner_id": widget.OwnerID(),
		"organisation_id": widget.OrganisationID(), "name": widget.Name(),
		"created_at": widget.CreatedAt(),
	})
	if err != nil {
		return fmt.Errorf("insert widget: %w", err)
	}
	return nil
}

func (r *WidgetRepository) ByID(ctx context.Context, id string) (domain.Widget, error) {
	row, err := r.rows.Get(ctx, id)
	if err != nil {
		return domain.Widget{}, fmt.Errorf("select widget: %w", err)
	}
	return domain.NewWidget(row.ID, row.OwnerID, row.OrganisationID, row.Name, row.CreatedAt)
}
