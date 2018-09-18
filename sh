protoc --descriptor_set_out=./test.protoset --include_imports -I. test.proto

protoc --descriptor_set_out=./test.protoset --include_imports -I. test.proto

protoc -I mdb/ mdb/tx.proto --go_out=plugins=mdb/feed
