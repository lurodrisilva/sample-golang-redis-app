package persistence_test

import (
	"context"
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/lucianosilva/sample-redis-app/internal/adapter/outbound/persistence"
	"github.com/lucianosilva/sample-redis-app/internal/domain/item"
)

func setupRepo(t *testing.T) (*persistence.RedisItemRepo, *miniredis.Miniredis) {
	t.Helper()
	mr := miniredis.RunT(t)
	client := redis.NewClient(&redis.Options{Addr: mr.Addr()})
	t.Cleanup(func() { client.Close() })
	return persistence.NewRedisItemRepo(client), mr
}

func TestRedisItemRepo_SaveAndFindByID(t *testing.T) {
	repo, _ := setupRepo(t)
	ctx := context.Background()

	i, err := item.New("Widget", "A test widget")
	require.NoError(t, err)

	err = repo.Save(ctx, i)
	require.NoError(t, err)

	found, err := repo.FindByID(ctx, i.ID())
	require.NoError(t, err)

	assert.Equal(t, i.ID().String(), found.ID().String())
	assert.Equal(t, i.Name(), found.Name())
	assert.Equal(t, i.Description(), found.Description())
}

func TestRedisItemRepo_FindByID_NotFound(t *testing.T) {
	repo, _ := setupRepo(t)
	ctx := context.Background()

	id := item.NewItemID()
	_, err := repo.FindByID(ctx, id)

	assert.True(t, err == item.ErrNotFound, "expected ErrNotFound, got: %v", err)
}

func TestRedisItemRepo_FindByID_RedisError(t *testing.T) {
	repo, mr := setupRepo(t)
	ctx := context.Background()

	i, err := item.New("Widget", "desc")
	require.NoError(t, err)
	require.NoError(t, repo.Save(ctx, i))

	mr.Close()

	_, err = repo.FindByID(ctx, i.ID())
	assert.Error(t, err)
}

func TestRedisItemRepo_Save_RedisError(t *testing.T) {
	repo, mr := setupRepo(t)
	ctx := context.Background()

	mr.Close()

	i, err := item.New("Widget", "desc")
	require.NoError(t, err)

	err = repo.Save(ctx, i)
	assert.Error(t, err)
}

func TestRedisItemRepo_FindByID_CorruptedJSON(t *testing.T) {
	repo, mr := setupRepo(t)
	ctx := context.Background()

	id := item.NewItemID()
	mr.Set("item:"+id.String(), "not-valid-json")

	_, err := repo.FindByID(ctx, id)
	assert.Error(t, err)
}

func TestRedisItemRepo_FindByID_InvalidItemID(t *testing.T) {
	repo, mr := setupRepo(t)
	ctx := context.Background()

	id := item.NewItemID()
	mr.Set("item:"+id.String(), `{"id":"not-a-uuid","name":"Widget","description":"desc","created_at":"2024-01-01T00:00:00Z"}`)

	_, err := repo.FindByID(ctx, id)
	assert.Error(t, err)
}

func TestRedisItemRepo_FindByID_InvalidCreatedAt(t *testing.T) {
	repo, mr := setupRepo(t)
	ctx := context.Background()

	id := item.NewItemID()
	mr.Set("item:"+id.String(), `{"id":"`+id.String()+`","name":"Widget","description":"desc","created_at":"not-a-date"}`)

	_, err := repo.FindByID(ctx, id)
	assert.Error(t, err)
}
