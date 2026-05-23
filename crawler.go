package crawler

import "github.com/liaozhai/set"

type Interface[T comparable] interface {
	Value() T
	Nodes() []T
}

type Transformer[T comparable] func(t T) Interface[T]

func run[T comparable](seed T, depth int, transform Transformer[T], st *set.Set[T]) []T {
	if depth <= 0 || st.Has(seed) {
		return nil
	}
	st.Add(seed)
	ch := make(chan T)
	go func(x T) {
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

	a := []T{}
	for s := range ch {
		a = append(a, s)
	}
	return a
}

func Crawl[T comparable](seed T, depth int, transform Transformer[T]) []T {
	return run(seed, depth, transform, set.New[T]())
}
