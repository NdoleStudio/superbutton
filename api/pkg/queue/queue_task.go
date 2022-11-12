package queue

// Task represents a push queue task
type Task struct {
	Method string
	URL    string
	Body   []byte
}
