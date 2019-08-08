.PHONY: clean

all: %.pb.go
	go build .

%.pb.go:
	go generate

clean:
	- rm *.pb.go
	go clean -cache

test: data_pb_test.go
	docker-compose -f sawtooth-default.yaml up -d
	sleep 5
	go test -timeout 300s -v -cover tp_test.go data_pb_test.go data_test.go

down:
	docker-compose -f sawtooth-default.yaml down