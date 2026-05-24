package crawler

import (
	"sync"

	"github.com/liaozhai/set"
)

// K - key type, V - value type
type Interface[K comparable, V any] interface {
	Value() V
	Nodes() []K
}

type Transformer[K comparable, V any] func(t K) Interface[K, V]

func run[K comparable, V any](seed K, depth int, transform Transformer[K, V], st *set.Set[K], wg *sync.WaitGroup, out chan V) {
	t := transform(seed)
	defer wg.Done()
	for _, n := range t.Nodes() {
		wg.Add(1)
		go run(n, depth-1, transform, st, wg, out)
	}
	out <- t.Value()
}

func Crawl[K comparable, V any](seed K, depth int, transform Transformer[K, V]) chan V {
	wg := &sync.WaitGroup{}
	wg.Add(1)
	out := make(chan V)
	go run(seed, depth, transform, set.New[K](), wg, out)
	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}
