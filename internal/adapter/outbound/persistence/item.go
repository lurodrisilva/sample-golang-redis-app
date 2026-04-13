package persistence

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/lucianosilva/sample-redis-app/internal/domain/item"
)

// Compile-time check: RedisItemRepo implements item.Repository.
var _ item.Repository = (*RedisItemRepo)(nil)

// RedisItemRepo implements item.Repository using Redis.
type RedisItemRepo struct {
	client *redis.Client
}

// NewRedisItemRepo creates a new RedisItemRepo.
func NewRedisItemRepo(client *redis.Client) *RedisItemRepo {
	return &RedisItemRepo{client: client}
}

type itemRecord struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	CreatedAt   string `json:"created_at"`
}

func (r *RedisItemRepo) Save(ctx context.Context, i *item.Item) error {
	rec := itemRecord{
		ID:          i.ID().String(),
		Name:        i.Name(),
		Description: i.Description(),
		CreatedAt:   i.CreatedAt().Format(time.RFC3339),
	}

	data, err := json.Marshal(rec)
	if err != nil {
		return fmt.Errorf("save item: marshal: %w", err)
	}

	key := redisKey(i.ID())
	if err := r.client.Set(ctx, key, data, 0).Err(); err != nil {
		return fmt.Errorf("save item: redis set: %w", err)
	}

	return nil
}

func (r *RedisItemRepo) FindByID(ctx context.Context, id item.ItemID) (*item.Item, error) {
	key := redisKey(id)

	data, err := r.client.Get(ctx, key).Bytes()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, item.ErrNotFound
		}
		return nil, fmt.Errorf("find item: redis get: %w", err)
	}

	var rec itemRecord
	if err := json.Unmarshal(data, &rec); err != nil {
		return nil, fmt.Errorf("find item: unmarshal: %w", err)
	}

	itemID, err := item.ParseItemID(rec.ID)
	if err != nil {
		return nil, fmt.Errorf("find item: %w", err)
	}

	createdAt, err := time.Parse(time.RFC3339, rec.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("find item: parse created_at: %w", err)
	}

	return item.Reconstitute(itemID, rec.Name, rec.Description, createdAt), nil
}

func redisKey(id item.ItemID) string {
	return "item:" + id.String()
}
