package abstractions

// ResponseHandler handler to implement when a request's response should be handled a specific way.
type ResponseHandler[T any] func(response T, errorMappings ErrorMappings) (interface{}, error)
