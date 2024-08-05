//go:build long

package main

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
)

type MockRepo struct {
	mock.Mock
}

func (m *MockRepo) CleanExpiredCarts(ctx context.Context, duration time.Duration) error {
	args := m.Called(ctx, duration)
	return args.Error(0)
}

func TestCleanExpiredCarts(t *testing.T) {
	mockRepo := new(MockRepo)
	mockRepo.On("CleanExpiredCarts", mock.Anything, time.Minute).Return(nil)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		ticker := time.NewTicker(time.Minute)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				mockRepo.CleanExpiredCarts(ctx, time.Minute)
			case <-ctx.Done():
				return
			}
		}
	}()

	time.Sleep(130 * time.Second)
	mockRepo.AssertNumberOfCalls(t, "CleanExpiredCarts", 2)
}
