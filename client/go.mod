module github.com/dairaga/sawtk/client

go 1.12

require (
	github.com/dairaga/log v0.0.0-20190611140521-2f471283f46f
	github.com/golang/protobuf v1.3.2
	github.com/hyperledger/sawtooth-sdk-go v0.1.2
)

replace github.com/hyperledger/sawtooth-sdk-go => ../../../hyperledger/sawtooth-sdk-go
