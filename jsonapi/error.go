package jsonapi

// Error represents a JSONAPI error value
type Error map[string]struct {
	Type   string `json:"type"`
	ID     string `json:"id"`
	Errors []struct {
		Status string
		Code   string
		Title  string
		Detail string
	}
}
