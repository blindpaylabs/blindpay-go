//go:build go1.23

package pagination

import "iter"

// All returns an iterator over the remaining items, for use with range
// (Go 1.23+). A failure is yielded once as the final (zero, err) pair.
//
//	for p, err := range pager.All() {
//		if err != nil {
//			return err
//		}
//		// ...
//	}
func (p *Pager[T]) All() iter.Seq2[T, error] {
	return func(yield func(T, error) bool) {
		for p.Next() {
			if !yield(p.Current(), nil) {
				return
			}
		}
		if err := p.Err(); err != nil {
			var zero T
			yield(zero, err)
		}
	}
}
