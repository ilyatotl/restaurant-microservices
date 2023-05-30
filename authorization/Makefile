.PHONY: generate
generate:
	rm -rf internal/app/pb/
	mkdir -p internal/app/pb
	protoc --proto_path=proto/ --go_out=internal/app/pb --go-grpc_out=internal/app/pb proto/*.proto