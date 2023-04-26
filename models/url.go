package models

import "time"

// URL holds a URL, how many times the URL has been submitted via the API. The time it was created and updated.
type URL struct {
	URL       string `query:"url"`
	Submitted int
	CreatedAt time.Time
	UpdatedAt time.Time
}
