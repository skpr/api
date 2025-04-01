#!/usr/bin/make -f

# Directory to store GRPC package.
DIR=./pb

# Image used to build the package.
IMAGE=skpr/grpc-go:latest

# Builds the client/server code.
build:
	# Cleaning up directories.
	rm -fR $(DIR)
	mkdir -p $(DIR)
	# Building image.
	docker build -t $(IMAGE) .
	# Building package.
	docker run -it -w $(PWD) -v $(PWD):$(PWD) $(IMAGE) /bin/bash -c 'protoc --go_out=./pb --go_opt=paths=source_relative --go-grpc_out=./pb --go-grpc_opt=paths=source_relative  *.proto'

docs:
	docker run --rm -v $(pwd)/doc:/out pseudomuto/protoc-gen-doc --doc_opt=markdown,docs.md

.PHONY: *
