package client

// BatchHeader is header of batch defined in sawtooth restful api.
type BatchHeader struct {
	SignerPublicKey string   `json:"signer_public_key,omitempty"`
	TransactionIds  []string `json:"transaction_ids,omitempty"`
}

// Batch ...
type Batch struct {
	Header          *BatchHeader  `json:"header,omitempty"`
	HeaderSignature string        `json:"header_signature,omitempty"`
	Transactions    []Transaction `json:"transactions,omitempty"`
}

// BatchList ...
type BatchList struct {
	Batches []Batch `json:"batches,omitempty"`
}

// BatchesResp response of batches.
type BatchesResp struct {
	Data   []Batch `json:"data,omitempty"`
	Head   string  `json:"head,omitempty"`
	Link   string  `json:"link,omitempty"`
	Paging *Paging `json:"paging,omitempty"`
}

// BatchResp response of batch.
type BatchResp struct {
	Data   *Batch `json:"data,omitempty"`
	Header string `json:"header,omitempty"`
	Link   string `json:"link,omitempty"`
}
