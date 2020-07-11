package repositories

import (
	"context"
	"time"

	"github.com/maykonlf/go-devkit/pkg/types/uuid"
	"github.com/maykonlf/webhook-middleware/internal/entities"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type RouteRepositoryI interface {
	ListRoutes(ctx context.Context, skip, take int64) (*[]entities.Route, error)
	AddRoute(ctx context.Context, route *entities.Route) (*uuid.UUID, error)
	GetRoute(ctx context.Context, id *uuid.UUID) (*entities.Route, error)
	UpdateRoute(ctx context.Context, id *uuid.UUID, route *entities.Route) error
	DeleteRoute(ctx context.Context, id *uuid.UUID) error
}

func NewRouteRepository(uri string, timeout time.Duration) RouteRepositoryI {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}

	return &RouteRepository{
		collection: client.Database("webhook").Collection("routes"),
		timeout:    timeout,
	}
}

type RouteRepository struct {
	collection *mongo.Collection
	timeout    time.Duration
}

func (r *RouteRepository) ListRoutes(ctx context.Context, skip, take int64) (*[]entities.Route, error) {
	cursor, err := r.collection.Find(ctx, bson.D{}, options.Find().SetSkip(skip).SetLimit(take))
	if err != nil {
		return nil, err
	}

	var routes []entities.Route
	for cursor.Next(context.Background()) {
		var route entities.Route
		if err := cursor.Decode(&route); err != nil {
			return nil, err
		}
		routes = append(routes, route)
	}

	return &routes, nil
}

func (r *RouteRepository) AddRoute(ctx context.Context, route *entities.Route) (*uuid.UUID, error) {
	route.ID = uuid.New()
	_, err := r.collection.InsertOne(ctx, route)
	return &route.ID, err
}

func (r *RouteRepository) GetRoute(ctx context.Context, id *uuid.UUID) (*entities.Route, error) {
	var route entities.Route

	err := r.collection.FindOne(ctx, bson.D{{"_id", *id}}).Decode(&route)
	return &route, err
}

func (r *RouteRepository) UpdateRoute(ctx context.Context, id *uuid.UUID, route *entities.Route) error {
	_, err := r.collection.UpdateOne(ctx, bson.D{{Key: "_id", Value: *id}}, bson.D{{Key: "$set", Value: route}})
	return err
}

func (r *RouteRepository) DeleteRoute(ctx context.Context, id *uuid.UUID) error {
	_, err := r.collection.DeleteOne(ctx, bson.D{{"_id", *id}})
	return err
}
