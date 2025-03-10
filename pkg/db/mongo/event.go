package mongo

import (
	"context"
	"errors"
	"fmt"
	"github.com/aaydin-tr/kyte"
	"github.com/chn555/property-service/pkg/property"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type EventState struct {
	client     *mongo.Client
	database   *mongo.Database
	collection *mongo.Collection
}

type EventStateConfig struct {
	DatabaseName   string
	CollectionName string
}

func NewEventState(client *mongo.Client, config *EventStateConfig) *EventState {
	database := client.Database(config.DatabaseName)
	collection := database.Collection(config.CollectionName)
	return &EventState{
		client:     client,
		database:   database,
		collection: collection,
	}
}

func (e *EventState) Close(ctx context.Context) error {
	return e.client.Disconnect(ctx)
}

func (e *EventState) SaveEvent(ctx context.Context, event *property.Event) error {
	_, err := e.collection.InsertOne(ctx, event)
	if err != nil {
		return err
	}
	return nil
}

func (e *EventState) GetEventsForFilter(ctx context.Context, filter *property.EventFilter, limit int, offset int) ([]*property.Event, error) {
	if filter == nil {
		return nil, fmt.Errorf("filter is nil")
	}

	mongoFilter, err := buildFilter(filter)
	if err != nil {
		return nil, fmt.Errorf("build filter: %w", err)
	}

	var events []*property.Event
	opts := options.Find().SetLimit(int64(limit)).SetSkip(int64(offset))
	cursor, err := e.collection.Find(ctx, mongoFilter, opts)
	if err != nil {
		return nil, fmt.Errorf("find: %w", err)
	}
	err = cursor.All(ctx, &events)
	if err != nil {
		return nil, fmt.Errorf("cursor all: %w", err)
	}

	return events, nil
}

func buildFilter(filter *property.EventFilter) (bson.D, error) {
	event := &property.Event{}
	nonEmptyFilter := false
	filterBuilder := kyte.Filter(kyte.Source(event))
	if filter.PropertyID != "" {
		filterBuilder = filterBuilder.Equal(&event.PropertyID, filter.PropertyID)
		nonEmptyFilter = true
	}
	if !filter.AfterTime.IsZero() {
		filterBuilder = filterBuilder.GreaterThanOrEqual(&event.Date, filter.AfterTime)
		nonEmptyFilter = true
	}
	if !filter.BeforeTime.IsZero() {
		filterBuilder = filterBuilder.LessThanOrEqual(&event.Date, filter.BeforeTime)
		nonEmptyFilter = true
	}
	if filter.AmountType == property.Expense {
		filterBuilder = filterBuilder.LessThan(&event.EventAmount, 0)
		nonEmptyFilter = true
	}
	if filter.AmountType == property.Income {
		filterBuilder = filterBuilder.GreaterThan(&event.EventAmount, 0)
		nonEmptyFilter = true
	}

	if !nonEmptyFilter {
		return nil, fmt.Errorf("no filter criteria specified")
	}

	return filterBuilder.Build()
}

func (e *EventState) GetMostRecentEventForFilter(ctx context.Context, filter *property.EventFilter) (*property.Event, bool, error) {
	if filter == nil {
		return nil, false, fmt.Errorf("filter is nil")
	}

	mongoFilter, err := buildFilter(filter)
	if err != nil {
		return nil, false, fmt.Errorf("build filter: %w", err)
	}

	event := &property.Event{}
	opts := options.FindOne().SetSort(bson.M{"date": -1})
	err = e.collection.FindOne(ctx, mongoFilter, opts).Decode(event)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, false, nil
		}
		return nil, false, fmt.Errorf("find: %w", err)
	}

	return event, true, nil
}
