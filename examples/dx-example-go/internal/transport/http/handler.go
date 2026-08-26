package http

import (
	"context"
	"time"

	httpx "github.com/datakaveri/dx-common-go/platform/http"

	"github.com/datakaveri/dx-example-go/internal/application"
	"github.com/datakaveri/dx-example-go/internal/domain"
)

type Handler struct{ app *application.Service }

func NewHandler(app *application.Service) *Handler { return &Handler{app: app} }

type CreateRequest struct {
	httpx.Actor
	Name          string `json:"name" validate:"required,min=1,max=120"`
	CorrelationID string `header:"X-Request-ID"`
}

type GetRequest struct {
	httpx.Actor
	ID string `path:"id" validate:"required,uuid"`
}

type WidgetResponse struct {
	ID             string    `json:"id"`
	OwnerID        string    `json:"ownerId"`
	OrganisationID string    `json:"organisationId"`
	Name           string    `json:"name"`
	CreatedAt      time.Time `json:"createdAt"`
}

func (h *Handler) Create(ctx context.Context, req CreateRequest) (httpx.Created[WidgetResponse], error) {
	widget, err := h.app.Create(ctx, application.CreateCommand{
		Subject: req.Subject, Name: req.Name, CorrelationID: req.CorrelationID,
	})
	if err != nil {
		return httpx.Created[WidgetResponse]{}, err
	}
	return httpx.Created[WidgetResponse]{Value: response(widget)}, nil
}

func (h *Handler) Get(ctx context.Context, req GetRequest) (WidgetResponse, error) {
	widget, err := h.app.Get(ctx, req.Subject, req.ID)
	if err != nil {
		return WidgetResponse{}, err
	}
	return response(widget), nil
}

func response(widget domain.Widget) WidgetResponse {
	return WidgetResponse{
		ID: widget.ID(), OwnerID: widget.OwnerID(), OrganisationID: widget.OrganisationID(),
		Name: widget.Name(), CreatedAt: widget.CreatedAt(),
	}
}
