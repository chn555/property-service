package property

import (
	"context"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/samber/lo"
	"reflect"
	"testing"
	"time"
)

func TestHandler_SaveEvent(t *testing.T) {
	type fields struct {
		store EventStore
	}
	type args struct {
		ctx        context.Context
		PropertyID string
		amount     float64
		date       time.Time
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    float64
		wantErr bool
	}{
		{
			name:   "no property ID",
			fields: fields{store: NewMockEventStore(nil, false)},
			args: args{
				ctx:        context.TODO(),
				PropertyID: "",
				amount:     100,
				date:       time.Now(),
			},
			want:    0,
			wantErr: true,
		},
		{
			name:   "zero amount",
			fields: fields{store: NewMockEventStore(nil, false)},
			args: args{
				ctx:        context.TODO(),
				PropertyID: gofakeit.Address().Address,
				amount:     0,
				date:       time.Now(),
			},
			want:    0,
			wantErr: false,
		},
		{
			name:   "zero date",
			fields: fields{store: NewMockEventStore(nil, false)},
			args: args{
				ctx:        context.TODO(),
				PropertyID: gofakeit.Address().Address,
				amount:     100,
				date:       time.Time{},
			},
			want:    0,
			wantErr: true,
		},
		{
			name:   "state error",
			fields: fields{store: NewMockEventStore(nil, true)},
			args: args{
				ctx:        context.TODO(),
				PropertyID: gofakeit.Address().Address,
				amount:     100,
				date:       time.Now(),
			},
			want:    0,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &Handler{
				store: tt.fields.store,
			}
			got, err := h.SaveEvent(tt.args.ctx, tt.args.PropertyID, tt.args.amount, tt.args.date)
			if (err != nil) != tt.wantErr {
				t.Errorf("SaveEvent() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("SaveEvent() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHandler_SaveEvent_ValidInput(t *testing.T) {
	propertyID := gofakeit.Address().Address
	oldBalance := gofakeit.Float64()
	// Create a mock store with successful event saving
	store := NewMockEventStore([]*Event{
		{
			PropertyID:       propertyID,
			PostEventBalance: oldBalance,
		},
	}, false)

	// Create a valid property ID

	// Set up test amount and date
	amount := gofakeit.Float64()
	date := time.Now()

	// Create handler with mock store
	h := &Handler{
		store: store,
	}

	// Call SaveEvent with valid inputs
	balance, err := h.SaveEvent(context.TODO(), propertyID, amount, date)

	// Verify no error occurred
	if err != nil {
		t.Errorf("SaveEvent() unexpected error = %v", err)
		return
	}
	// Verify balance is correct
	if balance != amount+oldBalance {
		t.Errorf("SaveEvent() got = %v, want %v", balance, amount+oldBalance)
	}
}

func TestHandler_GetPropertyEvents(t *testing.T) {
	type fields struct {
		store EventStore
	}
	type args struct {
		ctx        context.Context
		PropertyID string
		dateFrom   time.Time
		dateTo     time.Time
		sortOrder  SortOrder
		amountType AmountType
		offset     int
		limit      int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []*Event
		wantErr bool
	}{
		{
			name:   "no property ID",
			fields: fields{store: NewMockEventStore(nil, false)},
			args: args{
				ctx:        context.TODO(),
				PropertyID: "",
				dateFrom:   time.Now(),
				dateTo:     time.Now().Add(2 * time.Minute),
				sortOrder:  Descending,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name:   "dates have no time between them",
			fields: fields{store: NewMockEventStore(nil, false)},
			args: args{
				ctx:        context.TODO(),
				PropertyID: gofakeit.Address().Address,
				dateFrom:   time.Now(),
				dateTo:     time.Now().Add(-2 * time.Minute),
				sortOrder:  Descending,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name:   "state error",
			fields: fields{store: NewMockEventStore(nil, true)},
			args: args{
				ctx:        context.TODO(),
				PropertyID: gofakeit.Address().Address,
				dateFrom:   time.Now(),
				dateTo:     time.Now().Add(2 * time.Minute),
				sortOrder:  Descending,
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &Handler{
				store: tt.fields.store,
			}
			got, err := h.GetPropertyEvents(tt.args.ctx, tt.args.PropertyID, tt.args.dateFrom, tt.args.dateTo, tt.args.sortOrder, tt.args.amountType, tt.args.offset, tt.args.limit)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetPropertyEvents() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetPropertyEvents() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHandler_GetPropertyEvents_ValidInput_Ascending(t *testing.T) {
	seed, h := seedTestHandler()

	// Set up test parameters
	dateFrom := time.Time{}
	dateTo := time.Time{}
	sortOrder := Ascending

	// Call GetPropertyEvents with valid inputs
	gotEvents, err := h.GetPropertyEvents(context.TODO(), seed.propertyID, dateFrom, dateTo, sortOrder, All, 0, 0)

	// Verify no error occurred
	if err != nil {
		t.Errorf("GetPropertyEvents() unexpected error = %v", err)
		return
	}

	// Verify correct number of events returned
	if len(gotEvents) != seed.eventCount {
		t.Errorf("GetPropertyEvents() got %d events, want %d", len(gotEvents), seed.eventCount)
		return
	}

	// Verify events are sorted in ascending order
	for i := 1; i < len(gotEvents); i++ {
		if gotEvents[i].Date.Before(gotEvents[i-1].Date) {
			t.Errorf("GetPropertyEvents() events not sorted correctly")
			return
		}
	}
}

func TestHandler_GetPropertyEvents_ValidInput_Descending(t *testing.T) {
	seed, h := seedTestHandler()

	// Set up test parameters
	dateFrom := time.Time{}
	dateTo := time.Time{}
	sortOrder := Descending

	// Call GetPropertyEvents with valid inputs
	gotEvents, err := h.GetPropertyEvents(context.TODO(), seed.propertyID, dateFrom, dateTo, sortOrder, All, 0, 0)

	// Verify no error occurred
	if err != nil {
		t.Errorf("GetPropertyEvents() unexpected error = %v", err)
		return
	}

	// Verify correct number of events returned
	if len(gotEvents) != seed.eventCount {
		t.Errorf("GetPropertyEvents() got %d events, want %d", len(gotEvents), seed.eventCount)
		return
	}

	// Verify events are sorted in ascending order
	for i := 1; i < len(gotEvents); i++ {
		if gotEvents[i].Date.After(gotEvents[i-1].Date) {
			t.Errorf("GetPropertyEvents() events not sorted correctly")
			return
		}
	}
}
func TestHandler_GetPropertyEvents_ValidInput_InvalidSortOrder(t *testing.T) {
	seed, h := seedTestHandler()

	// Set up test parameters
	dateFrom := time.Time{}
	dateTo := time.Time{}
	sortOrder := SortOrder(23)

	// Call GetPropertyEvents with valid inputs
	gotEvents, err := h.GetPropertyEvents(context.TODO(), seed.propertyID, dateFrom, dateTo, sortOrder, All, 0, 0)

	// Verify no error occurred
	if err != nil {
		t.Errorf("GetPropertyEvents() unexpected error = %v", err)
		return
	}

	// Verify correct number of events returned
	if len(gotEvents) != seed.eventCount {
		t.Errorf("GetPropertyEvents() got %d events, want %d", len(gotEvents), seed.eventCount)
		return
	}

	// Verify events are sorted in ascending order
	for i := 1; i < len(gotEvents); i++ {
		if gotEvents[i].Date.After(gotEvents[i-1].Date) {
			t.Errorf("GetPropertyEvents() events not sorted correctly")
			return
		}
	}
}

func seedTestHandler() (seedInfo, *Handler) {
	// Create a property ID for testing
	propertyID := gofakeit.Address().Address

	var eventArr [100]*Event
	gofakeit.Slice(&eventArr)
	events := eventArr[:]

	// Create a mock store with the test events
	store := NewMockEventStore(events, false)

	// Create handler with mock store
	h := &Handler{
		store: store,
	}

	return seedInfo{
		propertyID: propertyID,
		eventCount: len(events),
		expense: len(lo.Filter(events, func(item *Event, index int) bool {
			return item.EventAmount < 0
		})),
		income: len(lo.Filter(events, func(item *Event, index int) bool {
			return item.EventAmount > 0
		})),
	}, h
}

type seedInfo struct {
	propertyID string
	eventCount int
	expense    int
	income     int
}
