module github.com/dairaga/sawtk/tx

go 1.12

require (
	github.com/dairaga/sawtk/signing v0.0.0-20190808015225-fe9be36a371b
	github.com/golang/protobuf v1.3.2
	github.com/hyperledger/sawtooth-sdk-go v0.1.2
	github.com/satori/go.uuid v1.2.0
)

replace (
	github.com/dairaga/sawtk/signing => ../signing
	github.com/hyperledger/sawtooth-sdk-go => ../../../hyperledger/sawtooth-sdk-go
)
