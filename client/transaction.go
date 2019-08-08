package client

// TransactionHeader is header of transaction defined in sawtooth restful api.
type TransactionHeader struct {
	BatcherPublicKey string   `json:"batcher_public_key,omitempty"`
	Dependencies     []string `json:"dependencies,omitempty"`
	FamilyName       string   `json:"family_name,omitempty"`
	FamilyVersion    string   `json:"family_version,omitempty"`
	Inputs           []string `json:"inputs,omitempty"`
	Nonce            string   `json:"nonce,omitempty"`
	Outputs          []string `json:"outputs,omitempty"`
	PayloadSha512    string   `json:"payload_sha512,omitempty"`
	SignerPublicKey  string   `json:"signer_public_key,omitempty"`
}

// Transaction ...
type Transaction struct {
	Header          *TransactionHeader `json:"header,omitempty"`
	HeaderSignature string             `json:"header_signature,omitempty"`
	Payload         string             `json:"payload,omitempty"`
}

// TransactionsResp ...
type TransactionsResp struct {
	Data   []Transaction `json:"data,omitempty"`
	Head   string        `json:"head,omitempty"`
	Link   string        `json:"link,omitempty"`
	Paging *Paging       `json:"paging,omitempty"`
}

// TransactionResp ...
type TransactionResp struct {
	Data *Block `json:"data,omitempty"`
	Head string `json:"head,omitempty"`
	Link string `json:"link,omitempty"`
}
