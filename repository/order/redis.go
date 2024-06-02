package order

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/Vkanhan/go-redis-microservice/model"
	"github.com/redis/go-redis/v9"
)

// RedisRepo struct contains a Redis client.
type RedisRepo struct {
	Client *redis.Client
}

// OrderIDKey generates a Redis key for the given order ID.
func OrderIDKey(id uint64) string {
	return fmt.Sprintf("order:%d", id)
}

// Insert inserts an order into Redis.
func (r *RedisRepo) Insert(ctx context.Context, order model.Order) error {
	// Marshal the order struct into JSON format
	data, err := json.Marshal(order)
	if err != nil {
		return fmt.Errorf("failed to encode error: %w", err)
	}

	// Generate a Redis key for the order using the order ID
	key := OrderIDKey(order.OrderID)

	// Start a Redis transaction
	txn := r.Client.TxPipeline()

	// Insert the order into Redis with the generated key, using SetNX to ensure it doesn't overwrite an existing entry
	res := txn.SetNX(ctx, key, string(data), 0)
	if err := res.Err(); err != nil {
		txn.Discard()	//discard the transaction if there is an error
		return fmt.Errorf("failed to set: %w", err)
	}

	// Add the key to a set named "orders"
	if err := txn.SAdd(ctx, "orders", key).Err(); err != nil {
		txn.Discard()
		return fmt.Errorf("failed to add orders set: %w", err)
	}

	// Execute the transaction
	if _, err := txn.Exec(ctx); err != nil {
		return fmt.Errorf("failed to exec: %w", err)
	}

	return nil
}

// ErrorNotExist is an error returned when an order does not exist in Redis.
var ErrorNotExist = errors.New("order does not exikst")

// FindByID retrieves an order from Redis by its ID.
func (r *RedisRepo) FindByID(ctx context.Context, id uint64) (model.Order, error) {
	// Generate a Redis key for the order using the given ID
	key := OrderIDKey(id)

	// Retrieve the order data from Redis using the generated key
	value, err := r.Client.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		// Return a custom error if the order does not exist in Redis
		return model.Order{}, ErrorNotExist
	} else if err != nil {
		// Return an error if there is a problem retrieving the order data
		return model.Order{}, fmt.Errorf("get order: %w", err)
	}

	// Unmarshal the JSON order data into an order struct
	var order model.Order

	err = json.Unmarshal([]byte(value), &order)
	if err != nil {
		return model.Order{}, fmt.Errorf("failed to decode order json: %w", err)
	}
	return order, nil

}

// DeleteByID deletes an order from Redis by its ID.
func (r *RedisRepo) DeleteByID(ctx context.Context, id uint64) error {

	// Generate a Redis key for the order using the given ID
	key := OrderIDKey(id)

	// Start a Redis transaction
	txn := r.Client.TxPipeline()

	// Delete the order data from Redis
	err := txn.Del(ctx, key).Err()
	if errors.Is(err, redis.Nil) {
		txn.Discard()	// Discard the transaction if the order does not exist
		return ErrorNotExist
	} else if err != nil {
		txn.Discard()
		return fmt.Errorf("get order: %w", err)
	}

	// Remove the key from the set named "orders"
	if err := txn.SRem(ctx, "orders", key).Err(); err != nil {
		txn.Discard()
		return fmt.Errorf("failed to remove the orders set: %w", err)
	}

	// Execute the transaction
	if _, err := txn.Exec(ctx); err != nil {
		return fmt.Errorf("failed to exec: %w", err)
	}

	return nil
}

// UpdateByID updates an existing order in Redis.
func (r *RedisRepo) UpdateByID(ctx context.Context, order model.Order) error {

	// Marshal the order struct into JSON format
	data, err := json.Marshal(order)
	if err != nil {
		return fmt.Errorf("failed to encode order: %w", err)
	}

	key := OrderIDKey(order.OrderID)

	// Update the order data in Redis, using SetXX to ensure it only updates existing entries
	err = r.Client.SetXX(ctx, key, string(data), 0).Err()
	if errors.Is(err, redis.Nil) {
		return ErrorNotExist
	} else if err != nil {
		return fmt.Errorf("set order: %w", err)
	}
	return nil

}

// FindAllPage struct defines the size and offset for pagination.
type FindAllPage struct {
	Size   uint64	// Number of records to return
	Offset uint64	// Offset for pagination
}

// FindResult struct holds the orders and the cursor for pagination.
type FindResult struct {
	Orders []model.Order	//List of orders
	Cursor uint64	// Cursor for pagination
}

// FindAll retrieves all orders from Redis with pagination.
func (r *RedisRepo) FindAll(ctx context.Context, page FindAllPage) (FindResult, error) {

	// Use SSCAN to get the keys with pagination
	res := r.Client.SScan(ctx, "orders", page.Offset, "*", int64(page.Size))

	keys, cursor, err := res.Result()
	if err != nil {
		return FindResult{}, fmt.Errorf("failed to get order id: %w", err)
	}

	// Return an empty result if no keys are found
	if len(keys) == 0 {
		return FindResult{
			Orders: []model.Order{},
		}, nil
	}

	// Use MGET to retrieve all the orders using the keys
	xs, err := r.Client.MGet(ctx, keys...).Result()
	if err != nil {
		return FindResult{}, fmt.Errorf("failed to get orders: %w", err)
	}

	orders := make([]model.Order, len(xs))

	// Unmarshal each order JSON data into an order struct
	for i, x := range xs {
		x := x.(string)

		var order model.Order

		err  := json.Unmarshal([]byte(x), &order)
		if err != nil {
			return FindResult{}, fmt.Errorf("failed to decode order json: %w", err)
		}

		orders[i] = order
	}

	return FindResult{
		Orders: orders,
		Cursor: cursor,
	}, nil
}
