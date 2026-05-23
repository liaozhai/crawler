package crawler

import "github.com/liaozhai/set"

// K - key type, V - value type
type Interface[K comparable, V any] interface {
	Value() V
	Nodes() []K
}

type Transformer[K comparable, V any] func(t K) Interface[K, V]

func run[K comparable, V any](seed K, depth int, transform Transformer[K, V], st *set.Set[K]) []V {
	if depth <= 0 || st.Has(seed) {
		return nil
	}
	st.Add(seed)
	ch := make(chan V)
	go func(x K) {
		defer close(ch)
		t := transform(x)
		v := t.Value()
		ch <- v
		for _, v := range t.Nodes() {
			b := run(v, depth-1, transform, st)
			for _, s := range b {
				ch <- s
			}
		}
	}(seed)

	a := []V{}
	for s := range ch {
		a = append(a, s)
	}
	return a
}

func Crawl[K comparable, V any](seed K, depth int, transform Transformer[K, V]) []V {
	return run(seed, depth, transform, set.New[K]())
}
