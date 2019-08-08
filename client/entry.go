package client

// Entry ...
type Entry struct {
	Address string `json:"address,omitempty"`
	Data    string `json:"data,omitempty"`
}

// EntriesResp ...
type EntriesResp struct {
	Data   []Entry `json:"data,omitempty"`
	Head   string  `json:"head,omitempty"`
	Link   string  `json:"link,omitempty"`
	Paging *Paging `json:"paging,omitempty"`
}

// EntryResp ...
type EntryResp struct {
	Data string `json:"data,omitempty"`
	Head string `json:"head,omitempty"`
	Link string `json:"link,omitempty"`
}
