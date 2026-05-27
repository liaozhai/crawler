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

func run[K comparable, V any](key K, depth int, transform Transformer[K, V], st *sync.Map, wg *sync.WaitGroup, out chan Result[K, V]) {
	defer wg.Done()
	if _, ok := st.Load(key); ok {
		return
	}
	t, err := transform(key)
	if err != nil {
		return
	}
	st.Store(key, struct{}{})
	for _, n := range t.Keys() {
		wg.Add(1)
		go run(n, depth-1, transform, st, wg, out)
	}
	out <- Result[K, V]{key, t.Value()}
}

func Crawl[K comparable, V any](seed K, depth int, transform Transformer[K, V], out chan Result[K, V]) {
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go run(seed, depth, transform, &sync.Map{}, wg, out)
	go func() {
		wg.Wait()
		close(out)
	}()
}
