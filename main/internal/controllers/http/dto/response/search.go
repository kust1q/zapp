package response

// For docs
type SearchResult struct {
	Users  []User  `json:"users,omitempty"`
	Tweets []Tweet `json:"tweets,omitempty"`
}
