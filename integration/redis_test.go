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
	tclog "github.com/testcontainers/testcontainers-go/log"
	tcredis "github.com/testcontainers/testcontainers-go/modules/redis"
)

func setupRedis(t *testing.T) *redis.Client {
	t.Helper()
	ctx := context.Background()

	t.Log("starting redis container")
	container, err := tcredis.Run(ctx, "redis:7-alpine",
		testcontainers.WithLogger(tclog.TestLogger(t)),
	)
	require.NoError(t, err)
	t.Cleanup(func() {
		t.Log("terminating redis container")
		require.NoError(t, container.Terminate(ctx))
	})

	addr, err := container.Endpoint(ctx, "")
	require.NoError(t, err)
	t.Logf("redis container ready at %s", addr)

	client := redis.NewClient(&redis.Options{Addr: addr})
	t.Cleanup(func() { client.Close() })

	return client
}

// TestPubSubViaFromChan verifies that Redis pub/sub messages flow through a FromChan source.
func TestPubSubViaFromChan(t *testing.T) {
	ctx := context.Background()
	client := setupRedis(t)

	t.Log("subscribing to test-channel")
	sub := client.Subscribe(ctx, "test-channel")
	t.Cleanup(func() { sub.Close() })

	t.Log("waiting for subscription confirmation")
	_, err := sub.Receive(ctx)
	require.NoError(t, err)
	t.Log("subscription confirmed")

	msgChan := make(chan string, 10)
	go func() {
		for msg := range sub.Channel() {
			t.Logf("received pub/sub message: %s", msg.Payload)
			msgChan <- msg.Payload
		}
	}()

	source := reactive.FromChan(msgChan)
	var mu sync.Mutex
	results := make([]string, 0)
	source.Observe(func(s string) error {
		mu.Lock()
		defer mu.Unlock()
		t.Logf("pipe observed: %s", s)
		results = append(results, s)
		return nil
	})
	source.Start()

	t.Log("publishing messages")
	require.NoError(t, client.Publish(ctx, "test-channel", "hello").Err())
	require.NoError(t, client.Publish(ctx, "test-channel", "world").Err())

	t.Log("waiting for results")
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

	t.Log("populating redis stream")
	for _, val := range []string{"first", "second", "third"} {
		require.NoError(t, client.XAdd(ctx, &redis.XAddArgs{
			Stream: "test-stream",
			Values: map[string]interface{}{"data": val},
		}).Err())
	}
	t.Log("stream populated")

	lastID := "0"
	itemsYielded := 0
	source := reactive.FromGeneratorWithDefaultBackoff(func() (*string, error) {
		// Signal done after all items have been yielded, before making another blocking XRead.
		if itemsYielded >= 3 {
			return nil, &reactive.GeneratorFinished{}
		}
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
		itemsYielded++
		t.Logf("generator yielded: %s", val)
		return &val, nil
	})

	results := make([]string, 0)
	source.Observe(func(s string) error {
		t.Logf("pipe observed: %s", s)
		results = append(results, s)
		return nil
	})
	source.UponClose(func() {
		t.Log("source completed")
		assert.Equal(t, []string{"first", "second", "third"}, results)
	})
	source.Start()

	t.Log("awaiting completion")
	source.AwaitCompletion()
}

// TestMapPipelineWithPubSub verifies that a Map operator correctly transforms
// messages arriving from a Redis pub/sub channel.
func TestMapPipelineWithPubSub(t *testing.T) {
	ctx := context.Background()
	client := setupRedis(t)

	t.Log("subscribing to numbers channel")
	sub := client.Subscribe(ctx, "numbers")
	t.Cleanup(func() { sub.Close() })

	t.Log("waiting for subscription confirmation")
	_, err := sub.Receive(ctx)
	require.NoError(t, err)
	t.Log("subscription confirmed")

	msgChan := make(chan string, 10)
	go func() {
		for msg := range sub.Channel() {
			t.Logf("received pub/sub message: %s", msg.Payload)
			msgChan <- msg.Payload
		}
	}()

	strSource := reactive.FromChan(msgChan)
	intSource := reactive.Map(strSource, func(s string) (int, error) {
		var n int
		_, err := fmt.Sscan(s, &n)
		t.Logf("mapped %q -> %d", s, n*2)
		return n * 2, err
	})
	strSource.Start()

	var mu sync.Mutex
	results := make([]int, 0)
	intSource.Observe(func(n int) error {
		mu.Lock()
		defer mu.Unlock()
		t.Logf("pipe observed: %d", n)
		results = append(results, n)
		return nil
	})

	t.Log("publishing messages")
	require.NoError(t, client.Publish(ctx, "numbers", "1").Err())
	require.NoError(t, client.Publish(ctx, "numbers", "2").Err())
	require.NoError(t, client.Publish(ctx, "numbers", "3").Err())

	t.Log("waiting for results")
	assert.Eventually(t, func() bool {
		mu.Lock()
		defer mu.Unlock()
		return len(results) == 3
	}, 3*time.Second, 50*time.Millisecond)

	mu.Lock()
	defer mu.Unlock()
	assert.Equal(t, []int{2, 4, 6}, results)
}
