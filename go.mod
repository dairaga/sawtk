module github.com/dairaga/sawtk

go 1.12

require (
	github.com/btcsuite/btcd v0.0.0-20190807005414-4063feeff79a
	github.com/btcsuite/btcutil v0.0.0-20190425235716-9e5f4b9a998d
	github.com/dairaga/log v0.0.0-20190611140521-2f471283f46f
	github.com/golang/protobuf v1.3.2
	github.com/hyperledger/sawtooth-sdk-go v0.1.2
	github.com/pebbe/zmq4 v1.0.0
	github.com/satori/go.uuid v1.2.0
	gitlab.com/dataforceme/sawtk v0.0.0-20190424075300-ae601e5109ec
)

replace github.com/hyperledger/sawtooth-sdk-go => ../../hyperledger/sawtooth-sdk-go
