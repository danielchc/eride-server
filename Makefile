clean:
	rm -r ./pb
	
gen:
	mkdir ./pb/ && protoc --proto_path=proto proto/*.proto  --go_out=:pb --go-grpc_out=:pb

run:
	go run cmd/main.go


.PHONY: clean gen run 