all: clean gen run

clean:
	rm -rf ./pb
	
gen:
	mkdir ./pb/ && protoc --proto_path=proto proto/*.proto  --go_out=:pb --go-grpc_out=:pb

run:
	go run cmd/main.go

rmdb:
	rm -rf *.db

.PHONY: clean gen run all