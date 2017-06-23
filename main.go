package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func display(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		select {
		case <-time.After(500 * time.Millisecond):
			fmt.Printf("display at %v\n", time.Now())
		case <-ctx.Done():
			fmt.Printf("context cancelled in display\n")
			return
		}
	}
}

func subprogress(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()

	select {
	case <-time.After(3000 * time.Millisecond):
		fmt.Printf("subprocess done!\n")
	case <-ctx.Done():
		fmt.Printf("context cancelled in subprogress\n")
		return
	}
}

func progress(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()

	subCtx, cancel := context.WithTimeout(ctx, 2500*time.Millisecond)
	defer cancel()

	wg.Add(1)
	go subprogress(subCtx, wg)

	for i := 0; i < 10; i++ {
		select {
		case <-time.After(2000 * time.Millisecond):
			fmt.Printf("  progress %d\n", 10*i)
		case <-ctx.Done():
			time.Sleep(100 * time.Millisecond)
			fmt.Printf("context cancelled in progress\n")
			return
		}
	}
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	sigCh := make(chan os.Signal)
	signal.Notify(sigCh, syscall.SIGINT)

	var wg sync.WaitGroup

	go func() {
		for sig := range sigCh {
			fmt.Printf("signal received: %v\n", sig)
			cancel()
			return
		}
	}()

	wg.Add(1)
	go progress(ctx, &wg)

	wg.Add(1)
	go display(ctx, &wg)

	wg.Wait()
}
