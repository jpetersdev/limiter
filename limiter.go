package limiter

import (
	"context"
	"errors"
	"golang.org/x/sync/semaphore"
)

var (
	MaxSizeReachedError = errors.New("max size reached")
	MinSizeReachedError = errors.New("min size reached")
)

type Limiter struct {
	maxSize     int64
	currentSize int64
	available   int64
	sem         *semaphore.Weighted
}

func NewLimiter(ctx context.Context, MaxSize, CurrentSize int64) (limiter *Limiter, err error) {
	s := semaphore.NewWeighted(MaxSize)
	err = s.Acquire(ctx, MaxSize-CurrentSize)
	return &Limiter{
		maxSize: MaxSize,
		currentSize: CurrentSize,
		available: MaxSize-CurrentSize,
		sem: s,
	}, err
}

func (l *Limiter) Increment() error {
	if l.available == 0 {
		return MaxSizeReachedError
	}
	l.sem.Release(1)
	l.currentSize++
	l.available--
	return nil
}

func (l *Limiter) Decrement(ctx context.Context) error {
	if l.currentSize == 0 {
		return MinSizeReachedError
	}
	err := l.sem.Acquire(ctx, 1)
	if err != nil { return err }
	l.currentSize++
	l.available--
	return nil
}

func (l *Limiter) Acquire(ctx context.Context) error {
	return l.sem.Acquire(ctx, 1)
}

func (l *Limiter) TryAcquire() bool {
	return l.sem.TryAcquire(1)
}

func (l *Limiter) Release() {
	l.sem.Release(1)
}
