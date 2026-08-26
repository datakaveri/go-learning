package events

import (
	"context"

	platformevents "github.com/datakaveri/dx-common-go/platform/events"

	"github.com/datakaveri/dx-example-go/internal/domain"
)

type WidgetCreated struct {
	WidgetID       string `json:"widgetId"`
	OwnerID        string `json:"ownerId"`
	OrganisationID string `json:"organisationId"`
}

var widgetCreated = platformevents.NewTopic[WidgetCreated]("widget.created").V(1)

type Writer struct {
	outbox     *platformevents.Outbox
	dispatcher *platformevents.Dispatcher
}

func NewWriter(outbox *platformevents.Outbox, dispatcher *platformevents.Dispatcher) *Writer {
	return &Writer{outbox: outbox, dispatcher: dispatcher}
}

func (w *Writer) WidgetCreated(ctx context.Context, widget domain.Widget, correlationID string) error {
	return platformevents.Publish(ctx, w.outbox, widgetCreated, WidgetCreated{
		WidgetID: widget.ID(), OwnerID: widget.OwnerID(), OrganisationID: widget.OrganisationID(),
	}, platformevents.WithCorrelationID(correlationID), platformevents.WithOccurredAt(widget.CreatedAt()))
}

func (w *Writer) Kick() { w.dispatcher.Kick() }
