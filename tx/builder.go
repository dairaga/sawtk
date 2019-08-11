package tx

import (
	"encoding/hex"
	"fmt"

	"github.com/dairaga/sawtk/signing"
	"github.com/golang/protobuf/proto"
	"github.com/hyperledger/sawtooth-sdk-go/protobuf/batch_pb2"
	"github.com/hyperledger/sawtooth-sdk-go/protobuf/transaction_pb2"
)

// BatchBuilder is to build sawtooth batch.
type BatchBuilder struct {
	signer *signing.Signer
}

// NewBatchBuilder returns a sawtooth batch builder.
func NewBatchBuilder(signer *signing.Signer) *BatchBuilder {
	return &BatchBuilder{signer: signer}
}

func (b *BatchBuilder) String() string {
	return fmt.Sprintf("signer: %s", b.signer.GetPublicKey().AsHex())
}

// BuildHeader returns sawtooth batch header.
func (b *BatchBuilder) BuildHeader(txs ...*transaction_pb2.Transaction) *batch_pb2.BatchHeader {
	ids := make([]string, len(txs))

	for i, x := range txs {
		ids[i] = x.HeaderSignature
	}

	return &batch_pb2.BatchHeader{
		TransactionIds:  ids,
		SignerPublicKey: b.signer.GetPublicKey().AsHex(),
	}
}

// Build returns sawtooth batch.
func (b *BatchBuilder) Build(txs ...*transaction_pb2.Transaction) (*batch_pb2.Batch, error) {
	header := b.BuildHeader(txs...)
	headerBytes, err := proto.Marshal(header)
	if err != nil {
		return nil, err
	}

	signature := hex.EncodeToString(b.signer.Sign(headerBytes))

	return &batch_pb2.Batch{
		Header:          headerBytes,
		Transactions:    txs,
		HeaderSignature: signature,
	}, nil
}

// BuildList returns sawtooth batch list.
func (b *BatchBuilder) BuildList(txs ...*transaction_pb2.Transaction) (*batch_pb2.BatchList, error) {
	batch, err := b.Build(txs...)
	if err != nil {
		return nil, err
	}

	return BatchList(batch), nil

}

// BatchList returns sawtooth batch list.
func BatchList(bs ...*batch_pb2.Batch) *batch_pb2.BatchList {
	return &batch_pb2.BatchList{Batches: bs}
}

// ----------------------------------------------------------------------------

// TransactionBuilder is to build sawtooth transaction.
type TransactionBuilder struct {
	batchSignerPublicKey string
	signer               *signing.Signer
}

// NewBuilder returns a TransactionBuilder.
func NewBuilder(batchSignerPublicKey string, signer *signing.Signer) *TransactionBuilder {
	return &TransactionBuilder{batchSignerPublicKey: batchSignerPublicKey, signer: signer}
}

func (b *TransactionBuilder) String() string {
	return fmt.Sprintf("batch signer: %s, signer: %s", b.batchSignerPublicKey, b.signer.GetPublicKey().AsHex())
}

// BuildHeader returns a sawtooth transaction header.
func (b *TransactionBuilder) BuildHeader(data *Data, dependencies ...string) *transaction_pb2.TransactionHeader {
	return &transaction_pb2.TransactionHeader{
		SignerPublicKey:  b.signer.GetPublicKey().AsHex(),
		FamilyName:       data.family,
		FamilyVersion:    data.version,
		Inputs:           data.inputs,
		Outputs:          data.outputs,
		Dependencies:     dependencies,
		BatcherPublicKey: b.batchSignerPublicKey,
		Nonce:            Nonce(),
		PayloadSha512:    signing.SHA512(data.payload),
	}
}

// Build returns a sawtooth transaction.
func (b *TransactionBuilder) Build(data *Data, dependencies ...string) (*transaction_pb2.Transaction, error) {
	header := b.BuildHeader(data, dependencies...)
	headerBytes, err := proto.Marshal(header)
	if err != nil {
		return nil, err
	}

	signature := hex.EncodeToString(b.signer.Sign(headerBytes))

	return &transaction_pb2.Transaction{
		Header:          headerBytes,
		HeaderSignature: signature,
		Payload:         data.payload,
	}, nil
}

// BuildBatch returns a sawtooth batch.
func (b *TransactionBuilder) BuildBatch(batchBuilder *BatchBuilder, data ...*Data) (*batch_pb2.Batch, error) {
	size := len(data)
	var err error
	txs := make([]*transaction_pb2.Transaction, size, size)
	for i := 0; i < size; i++ {
		txs[i], err = b.Build(data[i])
		if err != nil {
			return nil, err
		}
	}

	return batchBuilder.Build(txs...)
}
