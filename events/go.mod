module github.com/dairaga/sawtk/events

go 1.12

require (
	github.com/dairaga/sawtk/errors v0.0.0-20190807160034-4d93ba2d4f30
	github.com/golang/protobuf v1.3.2
	github.com/hyperledger/sawtooth-sdk-go v0.1.2
)

replace (
	github.com/dairaga/sawtk/errors => ../errors
	github.com/hyperledger/sawtooth-sdk-go => ../../../hyperledger/sawtooth-sdk-go
)
