package crawler

import (
	"sync"

	"github.com/liaozhai/set"
)

// K - key type, V - value type
type Interface[K comparable, V any] interface {
	Value() V
	Keys() []K
}

type Transformer[K comparable, V any] func(t K) Interface[K, V]

type Result[K comparable, V any] struct {
	Key   K
	Value V
}

func run[K comparable, V any](key K, depth int, transform Transformer[K, V], st *set.Set[K], wg *sync.WaitGroup, out chan Result[K, V]) {
	t := transform(key)
	defer wg.Done()
	for _, n := range t.Keys() {
		wg.Add(1)
		go run(n, depth-1, transform, st, wg, out)
	}
	out <- Result[K, V]{key, t.Value()}
}

func Crawl[K comparable, V any](seed K, depth int, transform Transformer[K, V]) chan Result[K, V] {
	wg := &sync.WaitGroup{}
	wg.Add(1)
	out := make(chan Result[K, V])
	go run(seed, depth, transform, set.New[K](), wg, out)
	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}
