package abstractions

type Builder[T any] interface {
	Build() (T, error)
}
