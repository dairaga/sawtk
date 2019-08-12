package tx

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	"github.com/dairaga/sawtk/signing"
	"github.com/golang/protobuf/proto"
	"github.com/hyperledger/sawtooth-sdk-go/protobuf/batch_pb2"
	"github.com/hyperledger/sawtooth-sdk-go/protobuf/transaction_pb2"
	uuid "github.com/satori/go.uuid"
)

// Data wraps sawtooth transaction payload, inputs and outputs.
type Data struct {
	family  string
	version string
	payload []byte
	inputs  []string
	outputs []string
}

func (d *Data) String() string {
	b64 := base64.StdEncoding.EncodeToString(d.payload)
	in := strings.Join(d.inputs, `","`)
	out := strings.Join(d.outputs, `","`)
	return fmt.Sprintf(`{"family": "%s", "version": "%s", "payload": "%s", "inputs": ["%s"], "outputs": ["%s"]}`, d.family, d.version, b64, in, out)
}

// New returns Data.
func New(family, version string, pb proto.Message, in, out []string) (*Data, error) {
	payload, err := proto.Marshal(pb)
	if err != nil {
		return nil, err
	}

	return &Data{
		family:  family,
		version: version,
		payload: payload,
		inputs:  in,
		outputs: out,
	}, nil
}

// ----------------------------------------------------------------------------

// Payload returns payload in data.
func (d *Data) Payload() []byte {
	return d.payload
}

// TxHeader return a transaction header.
func (d *Data) TxHeader(txkey, batchkey string, dependencies ...string) *transaction_pb2.TransactionHeader {
	return &transaction_pb2.TransactionHeader{
		SignerPublicKey:  txkey,
		FamilyName:       d.family,
		FamilyVersion:    d.version,
		Inputs:           d.inputs,
		Outputs:          d.outputs,
		Dependencies:     dependencies,
		BatcherPublicKey: batchkey,
		Nonce:            Nonce(),
		PayloadSha512:    signing.SHA512(d.payload),
	}
}

// ToTx returns a transaction.
func (d *Data) ToTx(txb *Builder, dependencies ...string) (*transaction_pb2.Transaction, error) {
	return txb.Build(d, dependencies...)
}

// ToBatch returns a batch including one transaction.
// bb is a batch builder.
// txb is a transaction builder.
// dependencies are transactions that the transaction depends on.
func (d *Data) ToBatch(bb *BatchBuilder, txb *Builder, dependencies ...string) (*batch_pb2.Batch, error) {
	tx, err := d.ToTx(txb, dependencies...)
	if err != nil {
		return nil, err
	}

	return bb.Build(tx)
}

// ToBatches returns a batch list including one batch with one transaction.
// bb is a batch builder.
// txb is a transaction builder.
// dependencies are transactions that the transaction depends on.
func (d *Data) ToBatches(bb *BatchBuilder, txb *Builder, dependencies ...string) (*batch_pb2.BatchList, error) {
	tx, err := d.ToTx(txb, dependencies...)
	if err != nil {
		return nil, err
	}

	return bb.BuildList(tx)
}

// ----------------------------------------------------------------------------

// Nonce returns nonce ID.
func Nonce() string {
	bytes, err := time.Now().MarshalBinary()
	if err != nil {
		return uuid.NewV4().String()
	}

	return hex.EncodeToString(bytes)
}
