.PHONY: proto install

proto:
# https://stackoverflow.com/questions/66168350/import-google-api-annotations-proto-was-not-found-or-had-errors-how-do-i-add/69428167#69428167
# Generate gRPC stubs
	protoc -I . \
		-I ./proto/googleapis \
		--go_out . \
		--go_opt paths=source_relative \
		--go-grpc_out . \
		--go-grpc_opt paths=source_relative \
		--go-grpc_opt require_unimplemented_servers=false \
		--grpc-gateway_out . \
		--grpc-gateway_opt logtostderr=true \
		--grpc-gateway_opt paths=source_relative \
		--grpc-gateway_opt generate_unbound_methods=true \
		./proto/api/meteostation.proto

install:
	go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest
	go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@latest
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
