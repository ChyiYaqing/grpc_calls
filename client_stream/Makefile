.PHONY: compiler

PROTOC_GEN_GO := $(GOPATH)/bin/protoc-gen-go
PROTOC_GEN_GO_GRPC := $(GOPATH)/bin/protoc-gen-go-grpc
PROTOC := $(shell which protoc)

# If protoc isn't on the path, see it to a target that's never up to date, so the install command always runs.
ifeq ($(PROTOC),)
	PROTOC = must-rebuild
endif

# Figure out which machine we're running on.
UNAME := $(shell uname)

$(PROTOC):
# Run the right installation command for the operating system.
ifeq ($(UNAME), Darwin)
	brew install protobuf
endif
ifeq ($(UNAME), Linux)
	sudo apt-get install protobuf-compiler
endif

# If $GOPATH/bin/protoc-gen-go does not exist, we'll run this command to install it.
$(PROTOC_GEN_GO):
	go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28

# If $GOPATH/bin/protoc-gen-go-grpc does not exist, we'll run this command to install it.
$(PROTOC_GEN_GO_GRPC):
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2

books.pb.go: proto/books.proto | $(PROTOC_GEN_GO) $(PROTOC)
	protoc --go_out=. --go_opt=paths=source_relative proto/books.proto

books_grpc.pb.go: proto/books.proto | $(PROTOC_GEN_GO_GRPC) $(PROTOC)
	protoc --go-grpc_out=. --go-grpc_opt=paths=source_relative proto/books.proto

# This is a "phony" target - an alias for the above command, so "make compile" still works.
compile: books.pb.go books_grpc.pb.go