package abstractions

// Represents a request option.
type ResponseOption interface {
	// GetKey returns the key to store the current option under.
	GetKey() ResponseOptionKey
}

// ResponseOptionKey represents a key to store a request option under.
type ResponseOptionKey struct {
	// The unique key for the option.
	Key string
}
