package abstractions

type RequestBuilder[T any] interface {
	Builder[T]
}
