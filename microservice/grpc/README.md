# Generate Protobuf

$protoc --proto_path=. --go_out=plugins=grpc:. hello.proto