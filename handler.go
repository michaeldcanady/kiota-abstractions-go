package abstractions

type Handler[T any] interface {
	Handle(T) error
}
