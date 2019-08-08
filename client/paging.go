package client

// Link response of batches with method POST
type Link struct {
	Link string `json:"link,omitempty"`
}

// Paging ...
type Paging struct {
	Start        string `json:"start,omitempty"`
	Limit        int32  `json:"limit,omitempty"`
	NextPosition string `json:"next_position,omitempty"`
	Next         string `json:"next,omitempty"`
}
