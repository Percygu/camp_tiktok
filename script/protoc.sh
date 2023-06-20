cd ../pkg/proto/
protoc --go_out=. --go-grpc_out=. ./*.proto