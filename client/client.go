package client

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/hyperledger/sawtooth-sdk-go/protobuf/setting_pb2"

	"github.com/dairaga/log"
	"github.com/golang/protobuf/proto"
	"github.com/hyperledger/sawtooth-sdk-go/protobuf/batch_pb2"
)

const (
	jsonType = "application/json"
	pbType   = "application/octet-stream"
)

// ----------------------------------------------------------------------------

func marshalToString(a interface{}) string {
	tmp, err := json.Marshal(a)
	if err != nil {
		return err.Error()
	}
	return string(tmp)
}

func queryString(pairs map[string]string) string {
	if len(pairs) <= 0 {
		return ""
	}

	builder := new(strings.Builder)

	for k, v := range pairs {
		if v != "" {
			builder.WriteString(fmt.Sprintf("&%s=%s", k, v))
		} else {
			builder.WriteString(fmt.Sprintf("&%s", k))
		}
	}

	if builder.Len() > 0 {
		return "?" + builder.String()[1:]
	}

	return ""
}

// ----------------------------------------------------------------------------

// Client sawtooth restful api client
type Client struct {
	endpoint string
	ref      *http.Client
}

// New a sawtooth restful api client
func New(endpoint string, timeout time.Duration) *Client {
	return &Client{endpoint: endpoint, ref: &http.Client{Timeout: timeout}}
}

// ----------------------------------------------------------------------------

func (cli *Client) String() string {
	return fmt.Sprintf(`{"endpoint": "%s", "timeout": "%v"}`, cli.endpoint, cli.ref.Timeout)
}

func (cli *Client) do(method, url, contentType string, data []byte) (ret *response) {
	ret = &response{}

	var reader io.Reader
	if data != nil {
		reader = bytes.NewReader(data)
	}
	req, err := http.NewRequest(method, url, reader)
	if err != nil {
		ret.lastErr = err
		return
	}
	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}

	resp, err := cli.ref.Do(req)
	if err != nil {
		ret.lastErr = err
		return
	}
	ret.fill(resp)
	return
}

// get query url with GET method.
func (cli *Client) get(url string) *response {
	return cli.do("GET", url, "", nil)
}

// postJSON queries url with POST method and JSON data.
func (cli *Client) postJSON(url string, data interface{}) *response {
	dataBytes, err := json.Marshal(data)
	if err != nil {
		return &response{lastErr: err}
	}

	return cli.do("POST", url, jsonType, dataBytes)
}

// postPB queries url with POST method and binary data from Protobuf.
func (cli *Client) postPB(url string, data proto.Message) *response {
	dataBytes, err := proto.Marshal(data)
	if err != nil {
		return &response{lastErr: err}
	}

	return cli.do("POST", url, pbType, dataBytes)
}

// SubmitBatches sumit batches to sawtooth restful api server.
func (cli *Client) SubmitBatches(batches *batch_pb2.BatchList) (*Link, error) {
	url := fmt.Sprintf("%s/batches", cli.endpoint)
	resp := cli.postPB(url, batches)

	ret := new(Link)

	if err := resp.handle(http.StatusAccepted, ret); err != nil {
		return nil, err
	}

	return ret, nil
}

// BatchesWithURL get batches with url
func (cli *Client) BatchesWithURL(url string) (*BatchesResp, error) {
	resp := cli.get(url)

	ret := new(BatchesResp)
	if err := resp.handle(http.StatusOK, resp); err != nil {
		return nil, err
	}

	return ret, nil
}

func (cli *Client) dataQS(head, start string, limit int, reverse string) string {
	if limit <= 0 {
		limit = 1000
	}

	return queryString(map[string]string{
		"head":    head,
		"start":   start,
		"limit":   strconv.Itoa(limit),
		"reverse": reverse,
	})
}

// Batches get batches
func (cli *Client) Batches(head, start string, limit int, reverse string) (*BatchesResp, error) {
	q := cli.dataQS(head, start, limit, reverse)

	url := fmt.Sprintf("%s/batches%s", cli.endpoint, q)

	return cli.BatchesWithURL(url)
}

// Batch get a batch
func (cli *Client) Batch(id string) (*BatchResp, error) {
	id = strings.TrimSpace(id)

	if id == "" {
		return nil, errors.New("id is required")
	}

	url := fmt.Sprintf("%s/batches/%s", cli.endpoint, id)

	resp := cli.get(url)

	ret := new(BatchResp)
	if err := resp.handle(http.StatusOK, ret); err != nil {
		return nil, err
	}
	return ret, nil
}

// BatchStatusesWithURL get batch statuses with url
func (cli *Client) BatchStatusesWithURL(url string) (*BatchStatuses, error) {
	resp := cli.get(url)

	ret := new(BatchStatuses)
	if err := resp.handle(http.StatusOK, ret); err != nil {
		return nil, err
	}
	return ret, nil
}

// BatchStatuses return a batch status
func (cli *Client) BatchStatuses(wait int, ids ...string) (*BatchStatuses, error) {
	if len(ids) == 0 {
		return nil, errors.New("one id at least")
	}

	id := strings.Join(ids, ",")

	if id == "" {
		return nil, errors.New("id required")
	}

	url := fmt.Sprintf("%s/batch_statuses?id=%s", cli.endpoint, id)

	if wait > 0 {
		url += fmt.Sprintf("&wait=%d", wait)
	}

	return cli.BatchStatusesWithURL(url)
}

// SubmitBatchStatuses get statuses of some batches
func (cli *Client) SubmitBatchStatuses(wait int32, ids ...string) (*BatchStatuses, error) {

	url := fmt.Sprintf("%s/batch_statuses", cli.endpoint)

	if wait > 0 {
		url += fmt.Sprintf("?wait=%d", wait)
	}

	resp := cli.postJSON(fmt.Sprintf("%s/batch_statuses", cli.endpoint), ids)

	ret := new(BatchStatuses)

	if err := resp.handle(http.StatusOK, ret); err != nil {
		return nil, err
	}

	return ret, nil
}

// States return address states
func (cli *Client) States(head, address, start string, limit int, reverse string) (*EntriesResp, error) {
	q := cli.dataQS(head, start, limit, reverse)

	url := fmt.Sprintf("%s/state%s", cli.endpoint, q)

	resp := cli.get(url)

	ret := new(EntriesResp)

	if err := resp.handle(http.StatusOK, ret); err != nil {
		return nil, err
	}
	return ret, nil
}

// State get state of an address
func (cli *Client) State(address, head string) (*EntryResp, error) {
	if address == "" {
		return nil, errors.New("address is required")
	}

	q := ""
	if head != "" {
		q = fmt.Sprintf("?head=%s", head)
	}

	url := fmt.Sprintf("%s/state/%s%s", cli.endpoint, address, q)

	resp := cli.get(url)
	ret := new(EntryResp)

	if err := resp.handle(http.StatusOK, ret); err != nil {
		return nil, err
	}
	return ret, nil
}

// StatePB ...
func (cli *Client) StatePB(address string, pb proto.Message) error {
	entry, err := cli.State(address, "")
	if err != nil {
		return err
	}

	if pb != nil {
		dataBytes, err := base64.StdEncoding.DecodeString(entry.Data)
		if err != nil {
			return err
		}

		return proto.Unmarshal(dataBytes, pb)
	}
	return nil
}

// Blocks return blocks
func (cli *Client) Blocks(head, start string, limit int, reverse string) (*BlocksResp, error) {
	q := cli.dataQS(head, start, limit, reverse)

	url := fmt.Sprintf("%s/blocks%s", cli.endpoint, q)

	resp := cli.get(url)

	ret := new(BlocksResp)

	if err := resp.handle(http.StatusOK, ret); err != nil {
		return nil, err
	}

	return ret, nil

}

// Block return a block
func (cli *Client) Block(id string) (*BlockResp, error) {
	if id == "" {
		return nil, errors.New("id is required")
	}

	url := fmt.Sprintf("%s/blocks/%s", cli.endpoint, id)
	resp := cli.get(url)

	ret := new(BlockResp)
	if err := resp.handle(http.StatusOK, ret); err != nil {
		return nil, err
	}
	return ret, nil
}

// Transactions return transactions
func (cli *Client) Transactions(head, start string, limit int, reverse string) (*TransactionsResp, error) {
	q := cli.dataQS(head, start, limit, reverse)

	url := fmt.Sprintf("%s/transactions%s", cli.endpoint, q)

	resp := cli.get(url)

	ret := new(TransactionsResp)
	if err := resp.handle(http.StatusOK, ret); err != nil {
		return nil, err
	}
	return ret, nil
}

// Transaction return a transaction
func (cli *Client) Transaction(id string) (*TransactionResp, error) {
	if id == "" {
		return nil, errors.New("id is required")
	}

	url := fmt.Sprintf("%s/transactions/%s", cli.endpoint, id)
	resp := cli.get(url)

	ret := new(TransactionResp)

	if err := resp.handle(http.StatusOK, ret); err != nil {
		return nil, err
	}

	return ret, nil
}

// SubmitBatchesResult ...
func (cli *Client) SubmitBatchesResult(batches *batch_pb2.BatchList) (*BatchStatuses, error) {

	link, err := cli.SubmitBatches(batches)
	if err != nil {
		log.Debugf("submit: %v", err)
		return nil, err
	}

	tmp := fmt.Sprintf("%s&wait=%.0f", link.Link, cli.ref.Timeout.Seconds())
	return cli.BatchStatusesWithURL(tmp)
}

// Data return data in chain state
func (cli *Client) Data(addr string, data proto.Message) error {
	return cli.StatePB(addr, data)
}

// Setting returns value in setting_tp.
func (cli *Client) Setting(addr string) (string, error) {
	s := new(setting_pb2.Setting)
	if err := cli.StatePB(addr, s); err != nil {
		return "", err
	}

	return s.Entries[0].Value, nil
}
