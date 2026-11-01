// Package option provides the functional-options helper shared by the rest of
// the library, so every constructor is configured the same way.
package option

// Func mutates a value of type T during construction.
type Func[T any] func(*T)

// Apply runs opts against t in order and returns t.
func Apply[T any](t *T, opts ...Func[T]) *T {
	for _, opt := range opts {
		opt(t)
	}
	return t
}
