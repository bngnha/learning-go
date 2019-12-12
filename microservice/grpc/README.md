https://grpc.io/docs/quickstart/go/
https://micro.mu/docs/toolkit.html

# Generate Protobuf

$protoc --proto_path=. --go_out=plugins=grpc:. hello.proto