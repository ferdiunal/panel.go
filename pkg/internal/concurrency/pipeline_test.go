package concurrency

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"time"
)

func TestAutoWorkers_DefaultRange(t *testing.T) {
	workers := AutoWorkers(0)
	if workers < 1 {
		t.Fatalf("expected at least 1 worker, got %d", workers)
	}
	if workers > 16 {
		t.Fatalf("expected at most 16 workers, got %d", workers)
	}
}

func TestClampWorkers_RespectsTotal(t *testing.T) {
	workers := ClampWorkers(32, 3)
	if workers != 3 {
		t.Fatalf("expected workers to clamp to total=3, got %d", workers)
	}
}

func TestMapOrdered_PreservesOrder(t *testing.T) {
	items := []int{10, 20, 30, 40}
	results, err := MapOrdered(context.Background(), items, 2, true, func(_ context.Context, idx int, item int) (int, error) {
		if idx%2 == 0 {
			time.Sleep(10 * time.Millisecond)
		}
		return item + idx, nil
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := []int{10, 21, 32, 43}
	for i := range expected {
		if results[i] != expected[i] {
			t.Fatalf("result mismatch at %d: expected %d got %d", i, expected[i], results[i])
		}
	}
}

func TestMapOrdered_FailFast(t *testing.T) {
	items := []int{1, 2, 3, 4, 5}
	var processed int32

	_, err := MapOrdered(context.Background(), items, 2, true, func(_ context.Context, _ int, item int) (int, error) {
		atomic.AddInt32(&processed, 1)
		if item == 2 {
			return 0, errors.New("boom")
		}
		time.Sleep(20 * time.Millisecond)
		return item, nil
	})
	if err == nil {
		t.Fatalf("expected fail-fast error")
	}
	if processed >= int32(len(items)) {
		t.Fatalf("expected cancellation to stop some work, processed=%d", processed)
	}
}

func TestMapOrdered_NoFailFastReturnsFirstError(t *testing.T) {
	items := []int{1, 2, 3}
	_, err := MapOrdered(context.Background(), items, 2, false, func(_ context.Context, _ int, item int) (int, error) {
		if item == 2 {
			return 0, errors.New("boom")
		}
		return item, nil
	})
	if err == nil {
		t.Fatalf("expected error in non-fail-fast mode")
	}
}

func TestMapOrdered_RespectsWorkerLimit(t *testing.T) {
	items := make([]int, 30)
	for i := range items {
		items[i] = i
	}

	const workerLimit = 3
	var (
		active    int32
		maxActive int32
	)

	_, err := MapOrdered(context.Background(), items, workerLimit, true, func(_ context.Context, _ int, item int) (int, error) {
		current := atomic.AddInt32(&active, 1)
		for {
			prev := atomic.LoadInt32(&maxActive)
			if current <= prev || atomic.CompareAndSwapInt32(&maxActive, prev, current) {
				break
			}
		}

		time.Sleep(10 * time.Millisecond)
		atomic.AddInt32(&active, -1)
		return item, nil
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if maxActive > workerLimit {
		t.Fatalf("expected max active workers <= %d, got %d", workerLimit, maxActive)
	}
}
