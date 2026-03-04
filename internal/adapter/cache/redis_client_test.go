package cache_test

import (
	"context"
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/davidlramirez95/towcommand/internal/adapter/cache"
)

func TestNewRedisClient(t *testing.T) {
	mr := miniredis.RunT(t)

	tests := []struct {
		name     string
		opts     cache.Options
		wantAddr string
	}{
		{
			name: "default pool size",
			opts: cache.Options{
				Host: mr.Host(),
				Port: mr.Server().Addr().Port,
			},
			wantAddr: mr.Addr(),
		},
		{
			name: "custom pool size",
			opts: cache.Options{
				Host:     mr.Host(),
				Port:     mr.Server().Addr().Port,
				PoolSize: 20,
			},
			wantAddr: mr.Addr(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := cache.NewRedisClient(tt.opts)
			defer client.Close()
			assert.NotNil(t, client)
		})
	}
}

func TestHealthCheck(t *testing.T) {
	mr := miniredis.RunT(t)
	ctx := context.Background()

	tests := []struct {
		name    string
		setup   func()
		wantErr bool
	}{
		{
			name:    "healthy server",
			setup:   func() {},
			wantErr: false,
		},
		{
			name: "closed server",
			setup: func() {
				mr.Close()
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := cache.NewRedisClient(cache.Options{
				Host: mr.Host(),
				Port: mr.Server().Addr().Port,
			})
			defer client.Close()

			tt.setup()

			err := cache.HealthCheck(ctx, client)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
