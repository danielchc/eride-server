all: clean gen config run

clean:
	rm -rf ./pb
	
gen:
	mkdir ./pb/ && protoc --proto_path=proto proto/*.proto  --go_out=:pb --go-grpc_out=:pb


config:
	cp config.yaml config.yaml.example

run:
	go run cmd/main.go

rmdb:
	rm -rf *.db

.PHONY: clean gen run all config rmdb