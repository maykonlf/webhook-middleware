package repositories

import (
	"context"
	"github.com/maykonlf/go-devkit/pkg/types/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"

	"github.com/maykonlf/webhook-middleware/internal/entities"
)

func TestCreateRoute(t *testing.T) {
	repository := NewRouteRepository("mongodb://localhost:27017/", 5*time.Second)

	route := &entities.Route{
		Body: "my body",
		URI:  "https://mockapi.io/123",
		Headers: map[string]string{
			"Authorization": "Basic 123",
		},
	}

	ctx := context.Background()

	id, err := repository.CreateRoute(ctx, route)
	assert.Nil(t, err, "error should be nil")
	assert.NotNil(t, id, "id should not be nil")

	routeFetched, err := repository.GetRoute(ctx, id)
	assert.Nil(t, err, "error should be nil")
	assert.Equal(t, route.Body, routeFetched.Body)
	assert.Equal(t, route.URI, routeFetched.URI)
	assert.Equal(t, id.String(), routeFetched.ID.String())

	routeFetched.Body = "new body"
	err = repository.UpdateRoute(ctx, id, routeFetched)
	assert.Nil(t, err)

	routeFetched, err = repository.GetRoute(ctx, id)
	assert.Nil(t, err)
	assert.Equal(t, routeFetched.Body, "new body")

	err = repository.DeleteRoute(ctx, id)
	assert.Nil(t, err, "delete error should be nil")

	var ids []uuid.UUID
	for i := 0; i < 3; i++ {
		id, err := repository.CreateRoute(ctx, &entities.Route{
			Body: "some body",
		})
		assert.Nil(t, err)
		ids = append(ids, *id)
	}

	listResult, err := repository.ListRoutes(ctx, 0, 10)
	assert.Nil(t, err)
	assert.Equal(t, 3, len(*listResult))

	for i := 0; i < len(*listResult); i++ {
		err := repository.DeleteRoute(ctx, &ids[i])
		assert.Nil(t, err)
	}
}
