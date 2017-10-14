package item

// Comment represents a comment
type Comment struct {
	ID           string `jsonapi:"primary,comments"`
	Body         string `jsonapi:"attr,body"`
	DiscussionID string `jsonapi:"attr,discussion_id"`
}
