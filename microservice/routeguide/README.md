# Generate Protobuf

protoc --proto_path=./routeguide/proto --go_out=plugins=grpc:./routeguide/proto ./routeguide/proto/route_guide.proto