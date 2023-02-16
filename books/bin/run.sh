#!/bin/bash

protoc --go_out=. --go_opt module=books --go-grpc_out=. --go-grpc_opt module=books -I . proto/books.proto
#protoc --go_out=. --go_opt module=books --go-grpc_out=. --go-grpc_opt module=books -I . proto/books.proto