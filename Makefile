BRANCH := $(shell git rev-parse --abbrev-ref HEAD)
COMMIT := $(shell git log -1 --format='%H')
DOCKER := $(shell which docker)
PACKAGES=$(shell go list ./...)
MOCKS_DIR = $(CURDIR)/tests/mocks

# don't override user values
ifeq (,$(VERSION))
  VERSION := $(shell git describe --tags 2>/dev/null)
  # if VERSION is empty, then populate it with branch's name and raw commit hash
  ifeq (,$(VERSION))
    VERSION := $(BRANCH)-$(COMMIT)
  endif
endif

# Update the ldflags with the app, client & server names
ldflags = -X github.com/cosmos/cosmos-sdk/version.Name=axmd \
	-X github.com/cosmos/cosmos-sdk/version.AppName=axmd \
	-X github.com/cosmos/cosmos-sdk/version.Version=$(VERSION) \
	-X github.com/cosmos/cosmos-sdk/version.Commit=$(COMMIT)

BUILD_FLAGS := -ldflags '$(ldflags)'

###########
# Install #
###########
all: install

install:
	@echo "--> ensure dependencies have not been modified"
	@go mod verify
	@echo "--> installing axmd"
	@go install $(BUILD_FLAGS) -mod=readonly ./cmd/axmd

init:
	./scripts/init-testchain.sh

##################
###  Building  ###
##################
build-all: proto
		"$(DOCKER)" run --rm -v "$(CURDIR):/axm-node" -w /axm-node golang:1.21-alpine sh ./scripts/build-all.sh

build: go.sum
		go build $(BUILD_FLAGS) ./cmd/axmd

go.sum: go.mod
		@echo "--> Ensure dependencies have not been modified"
		go mod verify
		go mod tidy

##################
###  Testing   ###
##################

$(MOCKS_DIR):
	mkdir -p $(MOCKS_DIR)

mocks: $(MOCKS_DIR)
	@go install github.com/golang/mock/mockgen@v1.6.0
	sh ./scripts/mockgen.sh

test:
	@go test -v -mod=readonly -tags=testing $(PACKAGES)

##################
###  Protobuf  ###
##################

protoVer=0.14.0
protoImageName=ghcr.io/cosmos/proto-builder:$(protoVer)
protoImage=$(DOCKER) run --rm  --network=host -v $(CURDIR):/workspace --workdir /workspace $(protoImageName)

proto-all: proto-clean proto-format proto-gen

proto-clean:
	find x -type f -iname *.pb.go -delete
	find x -type f -iname *.pb.gw.go -delete
	find util -type f -iname *.pb.go -delete

proto-gen:
	@echo "Generating protobuf files..."
	@$(protoImage) sh ./scripts/protocgen.sh
	@go mod tidy

proto-swagger-gen:
	@echo "Generating Protobuf Swagger"
	@$(protoImage) sh ./scripts/protoc-swagger-gen.sh

proto-format:
	@$(protoImage) find ./ -name "*.proto" -exec clang-format -i {} \;

proto-lint:
	@$(protoImage) buf lint proto/ --error-format=json

