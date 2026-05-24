package cache

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/require"
)

type fakeRedis struct {
	getVal string
	getErr error
	setErr error

	gotKey string
	gotVal any
	gotTTL time.Duration
}

func (f *fakeRedis) Get(_ context.Context, _ string) *redis.StringCmd {
	return redis.NewStringResult(f.getVal, f.getErr)
}

func (f *fakeRedis) Set(_ context.Context, key string, val any, ttl time.Duration) *redis.StatusCmd {
	f.gotKey, f.gotVal, f.gotTTL = key, val, ttl
	return redis.NewStatusResult("OK", f.setErr)
}

func TestGet_Hit(t *testing.T) {
	c := New(&fakeRedis{getVal: "https://example.com"}, time.Minute)
	got, err := c.Get(context.Background(), "abc")
	require.NoError(t, err)
	require.Equal(t, "https://example.com", got)
}

func TestGet_Miss(t *testing.T) {
	c := New(&fakeRedis{getErr: redis.Nil}, time.Minute)
	_, err := c.Get(context.Background(), "abc")
	require.ErrorIs(t, err, ErrNotFound)
}

func TestGet_Error(t *testing.T) {
	c := New(&fakeRedis{getErr: errors.New("boom")}, time.Minute)
	_, err := c.Get(context.Background(), "abc")
	require.ErrorContains(t, err, "boom")
}

func TestSet_PassesArgsThrough(t *testing.T) {
	f := &fakeRedis{}
	c := New(f, 5*time.Minute)
	err := c.Set(context.Background(), "abc", "https://example.com")
	require.NoError(t, err)
	require.Equal(t, "abc", f.gotKey)
	require.Equal(t, "https://example.com", f.gotVal)
	require.Equal(t, 5*time.Minute, f.gotTTL)
}

func TestSet_Error(t *testing.T) {
	c := New(&fakeRedis{setErr: errors.New("boom")}, time.Minute)
	err := c.Set(context.Background(), "abc", "x")
	require.ErrorContains(t, err, "boom")
}
