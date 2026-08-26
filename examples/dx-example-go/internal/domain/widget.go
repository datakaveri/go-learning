package domain

import (
	"strings"
	"time"

	platformerrors "github.com/datakaveri/dx-common-go/platform/errors"
)

// Widget is the teaching aggregate. Fields stay private so invalid state cannot
// bypass NewWidget or persistence reconstitution.
type Widget struct {
	id             string
	ownerID        string
	organisationID string
	name           string
	createdAt      time.Time
}

// NewWidget validates new state.
func NewWidget(id, ownerID, organisationID, name string, createdAt time.Time) (Widget, error) {
	name = strings.TrimSpace(name)
	switch {
	case id == "":
		return Widget{}, platformerrors.Validation("widget id is required")
	case ownerID == "":
		return Widget{}, platformerrors.Validation("widget owner is required")
	case organisationID == "":
		return Widget{}, platformerrors.Validation("widget organisation is required")
	case name == "":
		return Widget{}, platformerrors.Validation("widget name is required")
	case len(name) > 120:
		return Widget{}, platformerrors.Validation("widget name must not exceed 120 characters")
	case createdAt.IsZero():
		return Widget{}, platformerrors.Validation("widget creation time is required")
	}
	return Widget{
		id: id, ownerID: ownerID, organisationID: organisationID,
		name: name, createdAt: createdAt.UTC(),
	}, nil
}

func (w Widget) ID() string             { return w.id }
func (w Widget) OwnerID() string        { return w.ownerID }
func (w Widget) OrganisationID() string { return w.organisationID }
func (w Widget) Name() string           { return w.name }
func (w Widget) CreatedAt() time.Time   { return w.createdAt }
