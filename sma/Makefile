.PHONY: compiler

PROTOC_GEN_GO := $(GOPATH)/bin/protoc-gen-go
PROTOC_GEN_GO_GRPC := $(GOPATH)/bin/protoc-gen-go-grpc
PROTOC_GEN_GRPC_GATEWAY := $(GOPATH)/bin/protoc-gen-grpc-gateway
PROTOC_GEN_SWAGGER := $(GOPATH)/bin/protoc-gen-swagger 
PROTOC_GEN_OPENAPI := $(GOPATH)/bin/protoc-gen-openapiv2
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

$(PROTOC_GEN_GRPC_GATEWAY):
	go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest

$(PROTOC_GEN_SWAGGER):
	go install github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger@latest

$(PROTOC_GEN_OPENAPI):
	go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@latest


posts.pb.go: proto/posts.proto | $(PROTOC_GEN_GO) $(PROTOC)
	protoc -I ./proto --go_out=./proto --go_opt=paths=source_relative proto/posts.proto

posts_grpc.pb.go: proto/posts.proto | $(PROTOC_GEN_GO_GRPC) $(PROTOC)
	protoc -I ./proto --go-grpc_out=./proto --go-grpc_opt=paths=source_relative proto/posts.proto

posts.pb.gw.go: proto/posts.proto | $(PROTOC_GEN_GRPC_GATEWAY)
	protoc -I ./proto \
		--grpc-gateway_out=./proto \
		--grpc-gateway_opt logtostderr=true \
		--grpc-gateway_opt paths=source_relative \
		posts.proto

posts.swagger.json: proto/posts.proto | $(PROTOC_GEN)
	protoc -I ./proto \
		--openapiv2_out ./proto/openapiv2 \
		--openapiv2_opt logtostderr=true \
		posts.proto

# This is a "phony" target - an alias for the above command, so "make compile" still works.
compile: posts.pb.go posts_grpc.pb.go posts.pb.gw.go posts.swagger.json | $(PROTOC_GEN_GRPC_GATEWAY) $(PROTOC_GEN_SWAGGER) $(PROTOC_GEN_OPENAPI)