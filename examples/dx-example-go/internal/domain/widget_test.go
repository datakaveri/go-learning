package domain

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestNewWidget(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		id      string
		owner   string
		org     string
		label   string
		created time.Time
		wantErr bool
	}{
		{name: "valid", id: "w-1", owner: "u-1", org: "o-1", label: "Air quality", created: time.Now()},
		{name: "empty name", id: "w-1", owner: "u-1", org: "o-1", created: time.Now(), wantErr: true},
		{name: "missing organisation", id: "w-1", owner: "u-1", label: "Air quality", created: time.Now(), wantErr: true},
		{name: "missing time", id: "w-1", owner: "u-1", org: "o-1", label: "Air quality", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			_, err := NewWidget(tt.id, tt.owner, tt.org, tt.label, tt.created)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
		})
	}
}
