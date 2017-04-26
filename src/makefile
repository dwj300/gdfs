all: protobuf

protobuf:
	protoc -I contracts contracts/blob_server.proto --go_out=plugins=grpc:contracts
	protoc -I contracts contracts/metadata_server.proto --go_out=plugins=grpc:contracts
