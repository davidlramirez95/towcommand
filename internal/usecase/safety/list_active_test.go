package safetyuc

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/davidlramirez95/towcommand/internal/domain/safety"
)

func TestListActiveSOSUseCase_Execute(t *testing.T) {
	tests := []struct {
		name    string
		limit   int32
		alerts  []safety.SOSAlert
		lErr    error
		wantLen int
		wantLim int32
		wantErr bool
	}{
		{
			name:  "returns multiple active alerts",
			limit: 10,
			alerts: []safety.SOSAlert{
				{AlertID: "SOS-1", Resolved: false, Timestamp: time.Now()},
				{AlertID: "SOS-2", Resolved: false, Timestamp: time.Now()},
			},
			wantLen: 2,
			wantLim: 10,
		},
		{
			name:    "returns empty when no active alerts",
			limit:   10,
			alerts:  []safety.SOSAlert{},
			wantLen: 0,
			wantLim: 10,
		},
		{
			name:  "uses default limit when 0",
			limit: 0,
			alerts: []safety.SOSAlert{
				{AlertID: "SOS-1", Resolved: false},
			},
			wantLen: 1,
			wantLim: 50,
		},
		{
			name:  "uses default limit when negative",
			limit: -1,
			alerts: []safety.SOSAlert{
				{AlertID: "SOS-1", Resolved: false},
			},
			wantLen: 1,
			wantLim: 50,
		},
		{
			name:    "returns error from lister",
			limit:   10,
			lErr:    errors.New("db error"),
			wantErr: true,
			wantLim: 10,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lister := new(mockSOSActiveLister)

			lister.On("FindActive", mock.Anything, tt.wantLim).Return(tt.alerts, tt.lErr)

			uc := NewListActiveSOSUseCase(lister)
			got, err := uc.Execute(context.Background(), tt.limit)

			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Len(t, got, tt.wantLen)
			lister.AssertExpectations(t)
		})
	}
}
