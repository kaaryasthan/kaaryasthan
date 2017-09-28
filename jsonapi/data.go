package jsonapi

// Data represents a typical payload
type Data struct {
	Type       string            `json:"type"`
	ID         string            `json:"id"`
	Attributes map[string]string `json:"attributes"`
}
