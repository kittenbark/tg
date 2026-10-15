package tg

import (
	"context"
	"sync"
	"testing"
)

func TestDynamic(t *testing.T) {
	calledA, calledB := 0, 0
	a := func(ctx context.Context, upd *Update) error { calledA++; return nil }
	b := func(ctx context.Context, upd *Update) error { calledB++; return nil }

	handler, reload := Dynamic(a)

	if err := handler(context.Background(), nil); err != nil {
		t.Fatalf("handler: %v", err)
	}
	if calledA != 1 || calledB != 0 {
		t.Fatalf("expected initial handler called once, got calledA=%d calledB=%d", calledA, calledB)
	}

	reload(b)

	if err := handler(context.Background(), nil); err != nil {
		t.Fatalf("handler: %v", err)
	}
	if calledA != 1 || calledB != 1 {
		t.Fatalf("expected reloaded handler called once, got calledA=%d calledB=%d", calledA, calledB)
	}
}

func TestDynamicConcurrent(t *testing.T) {
	handler, reload := Dynamic(func(ctx context.Context, upd *Update) error { return nil })

	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		for range 1000 {
			_ = handler(context.Background(), nil)
		}
	}()
	go func() {
		defer wg.Done()
		for range 1000 {
			reload(func(ctx context.Context, upd *Update) error { return nil })
		}
	}()
	wg.Wait()
}
