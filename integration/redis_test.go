//go:build integration

package integration_test

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/AndreasChristianson/gopher-pipes/reactive"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	tcredis "github.com/testcontainers/testcontainers-go/modules/redis"
)

func setupRedis(t *testing.T) *redis.Client {
	t.Helper()
	ctx := context.Background()

	container, err := tcredis.Run(ctx, "redis:7-alpine",
		testcontainers.WithLogger(testcontainers.TestLogger(t)),
	)
	require.NoError(t, err)
	t.Cleanup(func() {
		require.NoError(t, container.Terminate(ctx))
	})

	addr, err := container.Endpoint(ctx, "")
	require.NoError(t, err)

	client := redis.NewClient(&redis.Options{Addr: addr})
	t.Cleanup(func() { client.Close() })

	return client
}

// TestPubSubViaFromChan verifies that Redis pub/sub messages flow through a FromChan source.
func TestPubSubViaFromChan(t *testing.T) {
	ctx := context.Background()
	client := setupRedis(t)

	sub := client.Subscribe(ctx, "test-channel")
	t.Cleanup(func() { sub.Close() })

	// Wait for subscription to be established before publishing.
	_, err := sub.Receive(ctx)
	require.NoError(t, err)

	msgChan := make(chan string, 10)
	go func() {
		for msg := range sub.Channel() {
			msgChan <- msg.Payload
		}
	}()

	source := reactive.FromChan(msgChan)
	var mu sync.Mutex
	results := make([]string, 0)
	source.Observe(func(s string) error {
		mu.Lock()
		defer mu.Unlock()
		results = append(results, s)
		return nil
	})
	source.Start()

	require.NoError(t, client.Publish(ctx, "test-channel", "hello").Err())
	require.NoError(t, client.Publish(ctx, "test-channel", "world").Err())

	assert.Eventually(t, func() bool {
		mu.Lock()
		defer mu.Unlock()
		return len(results) == 2
	}, 3*time.Second, 50*time.Millisecond)

	mu.Lock()
	defer mu.Unlock()
	assert.Equal(t, []string{"hello", "world"}, results)
}

// TestStreamPollingViaFromGenerator verifies that pre-populated Redis stream entries
// are consumed in order by a FromGeneratorWithDefaultBackoff source using XREAD.
func TestStreamPollingViaFromGenerator(t *testing.T) {
	ctx := context.Background()
	client := setupRedis(t)

	for _, val := range []string{"first", "second", "third"} {
		require.NoError(t, client.XAdd(ctx, &redis.XAddArgs{
			Stream: "test-stream",
			Values: map[string]interface{}{"data": val},
		}).Err())
	}

	lastID := "0"
	source := reactive.FromGeneratorWithDefaultBackoff(func() (*string, error) {
		res, err := client.XRead(ctx, &redis.XReadArgs{
			Streams: []string{"test-stream", lastID},
			Count:   1,
		}).Result()
		if err != nil || len(res) == 0 || len(res[0].Messages) == 0 {
			return nil, nil
		}
		msg := res[0].Messages[0]
		lastID = msg.ID
		val := fmt.Sprint(msg.Values["data"])
		return &val, nil
	})

	var mu sync.Mutex
	results := make([]string, 0)
	source.Observe(func(s string) error {
		mu.Lock()
		defer mu.Unlock()
		results = append(results, s)
		return nil
	})
	source.Start()

	assert.Eventually(t, func() bool {
		mu.Lock()
		defer mu.Unlock()
		return len(results) == 3
	}, 3*time.Second, 50*time.Millisecond)

	source.Cancel()
	source.AwaitCompletion()

	mu.Lock()
	defer mu.Unlock()
	assert.Equal(t, []string{"first", "second", "third"}, results)
}

// TestMapPipelineWithPubSub verifies that a Map operator correctly transforms
// messages arriving from a Redis pub/sub channel.
func TestMapPipelineWithPubSub(t *testing.T) {
	ctx := context.Background()
	client := setupRedis(t)

	sub := client.Subscribe(ctx, "numbers")
	t.Cleanup(func() { sub.Close() })

	_, err := sub.Receive(ctx)
	require.NoError(t, err)

	msgChan := make(chan string, 10)
	go func() {
		for msg := range sub.Channel() {
			msgChan <- msg.Payload
		}
	}()

	strSource := reactive.FromChan(msgChan)
	intSource := reactive.Map(strSource, func(s string) (int, error) {
		var n int
		_, err := fmt.Sscan(s, &n)
		return n * 2, err
	})

	var mu sync.Mutex
	results := make([]int, 0)
	intSource.Observe(func(n int) error {
		mu.Lock()
		defer mu.Unlock()
		results = append(results, n)
		return nil
	})

	require.NoError(t, client.Publish(ctx, "numbers", "1").Err())
	require.NoError(t, client.Publish(ctx, "numbers", "2").Err())
	require.NoError(t, client.Publish(ctx, "numbers", "3").Err())

	assert.Eventually(t, func() bool {
		mu.Lock()
		defer mu.Unlock()
		return len(results) == 3
	}, 3*time.Second, 50*time.Millisecond)

	mu.Lock()
	defer mu.Unlock()
	assert.Equal(t, []int{2, 4, 6}, results)
}
