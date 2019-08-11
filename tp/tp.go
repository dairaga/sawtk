package tp

import (
	"syscall"

	"github.com/hyperledger/sawtooth-sdk-go/processor"
)

// Run starts a Sawtooth Transaction Processor.
func Run(endpoint string, handler *Handler) error {

	processor := processor.NewTransactionProcessor(endpoint)
	processor.AddHandler(handler)
	processor.ShutdownOnSignal(syscall.SIGINT, syscall.SIGTERM)
	return processor.Start()
}

// Must starts a Sawtooth Transaction Processor.
// Panic if starting failure.
func Must(endpoint string, handler *Handler) {
	if err := Run(endpoint, handler); err != nil {
		panic(err)
	}
}
