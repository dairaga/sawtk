package client

import "encoding/json"

// Batch Status
const (
	BSCommitted = "COMMITTED"
	BSInvalid   = "INVALID"
	BSPending   = "PENDING"
	BSUnknown   = "UNKNOWN"
)

// InvalidTransaction describes invalid transaction defined in sawtooth restful api.
type InvalidTransaction struct {
	ID           string `json:"id,omitempty"`
	Message      string `json:"message,omitempty"`
	ExtendedData string `json:"extended_data,omitempty"`
}

// BatchStatus describes status of batch defined in sawtooth restful api.
type BatchStatus struct {
	ID                  string               `json:"id,omitempty"`
	Status              string               `json:"status,omitempty"`
	InvalidTransactions []InvalidTransaction `json:"invalid_transactions,omitempty"`
}

// IsCommitted ...
func (bs *BatchStatus) IsCommitted() bool {
	return bs.Status == BSCommitted
}

// IsPending ...
func (bs *BatchStatus) IsPending() bool {
	return bs.Status == BSPending
}

// IsInvalid ...
func (bs *BatchStatus) IsInvalid() bool {
	return bs.Status == BSInvalid
}

// IsUnknown ...
func (bs *BatchStatus) IsUnknown() bool {
	return bs.Status == BSUnknown
}

// BatchStatuses response of batch_statuses
type BatchStatuses struct {
	Data []BatchStatus `json:"data,omitempty"`
	Link string        `json:"link,omitempty"`
}

// IsOK ...
func (bss *BatchStatuses) IsOK() bool {
	for _, x := range bss.Data {
		if !x.IsCommitted() {
			return false
		}
	}

	return true
}

func (bss *BatchStatuses) filter(status string) []BatchStatus {
	var ret []BatchStatus
	for _, x := range bss.Data {
		if x.Status == status {
			ret = append(ret, x)
		}
	}
	return ret
}

// Committed returns committed batches.
func (bss *BatchStatuses) Committed() []BatchStatus {
	return bss.filter(BSCommitted)
}

// Invalid returns invalid batches.
func (bss *BatchStatuses) Invalid() []BatchStatus {
	return bss.filter(BSInvalid)
}

// Pending returns pending batches.
func (bss *BatchStatuses) Pending() []BatchStatus {
	return bss.filter(BSPending)
}

// Unknown returns unknown batches.
func (bss *BatchStatuses) Unknown() []BatchStatus {
	return bss.filter(BSUnknown)
}

func (bss *BatchStatuses) Error() string {
	var tmp []BatchStatus

	for _, x := range bss.Data {
		if !x.IsCommitted() {
			tmp = append(tmp, x)
		}
	}

	if len(tmp) > 0 {
		tmpbytes, err := json.Marshal(tmp)
		if err != nil {
			return err.Error()
		}

		return string(tmpbytes)
	}

	return ""
}

func (bss *BatchStatuses) String() string {
	tmpbytes, err := json.Marshal(bss)
	if err != nil {
		return err.Error()
	}

	return string(tmpbytes)
}
