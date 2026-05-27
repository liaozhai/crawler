package crawler

import (
	"sync"
)

// K - key type, V - value type
type Vertice[K comparable, V any] interface {
	Value() V
	Keys() []K
}

type Transformer[K comparable, V any] func(t K) (Vertice[K, V], error)

type Result[K comparable, V any] struct {
	Key   K
	Value V
}

type locker[K comparable] struct {
	mu    sync.Mutex
	cache map[K]struct{}
}

func newLocker[K comparable]() *locker[K] {
	return &locker[K]{sync.Mutex{}, make(map[K]struct{})}
}

func run[K comparable, V any](key K, depth int, transform Transformer[K, V], lck *locker[K], wg *sync.WaitGroup, out chan Result[K, V]) {
	defer wg.Done()
	lck.mu.Lock()
	if _, ok := lck.cache[key]; ok {
		lck.mu.Unlock()
		return
	}
	t, err := transform(key)
	if err != nil {
		lck.mu.Unlock()
		return
	}
	lck.cache[key] = struct{}{}
	out <- Result[K, V]{key, t.Value()}
	lck.mu.Unlock()
	for _, n := range t.Keys() {
		wg.Add(1)
		go run(n, depth-1, transform, lck, wg, out)
	}
}

func Crawl[K comparable, V any](seed K, depth int, transform Transformer[K, V], out chan Result[K, V]) {
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go run(seed, depth, transform, newLocker[K](), wg, out)
	go func() {
		wg.Wait()
		close(out)
	}()
}
