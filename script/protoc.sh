cd ../pkg/proto/
protoc --go-grpc_out=require_unimplemented_servers=false:.  --go_out=. ./*.proto