package client

// BlockHeader header of block defined in sawtooth restful api.
type BlockHeader struct {
	BlockNum        int32    `json:"block_num,omitempty"`
	PreviousBlockID string   `json:"previous_block_id,omitempty"`
	SignerPublicKey string   `json:"signer_public_key,omitempty"`
	BatchIds        []string `json:"batch_ids,omitempty"`
	Consensus       string   `json:"consensus,omitempty"`
	StateRootHash   string   `json:"state_root_hash,omitempty"`
}

// Block ...
type Block struct {
	Header          *BlockHeader `json:"header,omitempty"`
	HeaderSignature string       `json:"header_signature,omitempty"`
	Batches         []Batch      `json:"batches,omitempty"`
}

// BlocksResp ...
type BlocksResp struct {
	Data   []Block `json:"data,omitempty"`
	Head   string  `json:"head,omitempty"`
	Link   string  `json:"link,omitempty"`
	Paging *Paging `json:"paging,omitempty"`
}

// BlockResp ...
type BlockResp struct {
	Data *Block `json:"data,omitempty"`
	Head string `json:"head,omitempty"`
	Link string `json:"link,omitempty"`
}
