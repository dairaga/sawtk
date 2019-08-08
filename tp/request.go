package tp

//go:generate protoc -I . --go_out=plugins=grpc:../../../../ request.proto

import "github.com/golang/protobuf/proto"

// NewTPRequest returns a SawTK TPRequest.
func NewTPRequest(cmd int32, data proto.Message) (*TPRequest, error) {
	req := new(TPRequest)
	req.Cmd = cmd

	if data != nil {
		dataBytes, err := proto.Marshal(data)
		if err != nil {
			return nil, err
		}
		req.Payload = dataBytes
	}

	return req, nil
}

// ToBytes converts SawTK TPRequest to bytes.
func (r *TPRequest) ToBytes() ([]byte, error) {
	return proto.Marshal(r)
}

// NewTPRequestBytes returns bytes encoding from request.
func NewTPRequestBytes(cmd int32, data proto.Message) ([]byte, error) {
	req, err := NewTPRequest(cmd, data)
	if err != nil {
		return nil, err
	}

	return req.ToBytes()
}

// UnmarshalTPRequest returns data and command in TPReqest.
func UnmarshalTPRequest(data []byte, pb proto.Message) (int32, error) {
	req := new(TPRequest)
	if err := proto.Unmarshal(data, req); err != nil {
		return 0, err
	}

	if err := proto.Unmarshal(req.Payload, pb); err != nil {
		return 0, err
	}

	return req.Cmd, nil
}
