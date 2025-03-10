package property

import (
	"context"
	"github.com/brianvoe/gofakeit/v7"
	"testing"
	"time"
)

type MockEventStore struct {
	events []*Event
	err    bool
}

func (m MockEventStore) SaveEvent(ctx context.Context, event *Event) error {
	if m.err {
		return gofakeit.Error()
	}
	m.events = append(m.events, event)
	return nil
}

func (m MockEventStore) GetEventsForFilter(ctx context.Context, filter *EventFilter, limit int, offset int) ([]*Event, error) {
	if m.err {
		return nil, gofakeit.Error()
	}
	if limit > 0 {

		return m.events[:limit], nil
	}
	return m.events, nil
}

func (m MockEventStore) GetMostRecentEventForFilter(ctx context.Context, filter *EventFilter) (*Event, bool, error) {
	if m.err {
		return nil, false, gofakeit.Error()
	}
	if len(m.events) == 0 {
		return nil, false, nil
	}
	return m.events[0], true, nil
}

func NewMockEventStore(events []*Event, shouldThrowErr bool) *MockEventStore {
	return &MockEventStore{
		events: events,
		err:    shouldThrowErr,
	}
}

func TestHandler_GetBalance(t *testing.T) {
	type fields struct {
		store EventStore
	}
	type args struct {
		ctx        context.Context
		PropertyID string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    float64
		wantErr bool
	}{
		{
			name: "no property ID",
			fields: fields{
				store: NewMockEventStore(nil, false),
			},
			args: args{
				ctx: context.TODO(),
			},
			want:    0,
			wantErr: true,
		},
		{
			name: "no event for property ID",
			fields: fields{
				store: NewMockEventStore(nil, false),
			},
			args: args{
				ctx:        context.TODO(),
				PropertyID: gofakeit.Address().Address,
			},
			want:    0,
			wantErr: false,
		},
		{
			name: "err from state",
			fields: fields{
				store: NewMockEventStore(nil, true),
			},
			args: args{
				ctx:        context.TODO(),
				PropertyID: gofakeit.Address().Address,
			},
			want:    0,
			wantErr: true,
		},
		{
			name: "event exists for property ID",
			fields: fields{
				store: NewMockEventStore([]*Event{{
					PropertyID:       "propID",
					EventAmount:      142228,
					PostEventBalance: 142229,
					Date:             time.Now(),
				}}, false),
			},
			args: args{
				ctx:        context.TODO(),
				PropertyID: "propID",
			},
			want:    142229,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &Handler{
				store: tt.fields.store,
			}
			got, err := h.GetBalance(tt.args.ctx, tt.args.PropertyID)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetBalance() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetBalance() got = %v, want %v", got, tt.want)
			}
		})
	}
}
