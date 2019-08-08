package tx

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	"github.com/golang/protobuf/proto"
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

// Nonce returns nonce ID.
func Nonce() string {
	bytes, err := time.Now().MarshalBinary()
	if err != nil {
		return uuid.NewV4().String()
	}

	return hex.EncodeToString(bytes)
}
