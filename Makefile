GC=go build
BUILD_NODE_PAR = -ldflags "-X github.com/polynetwork/poly/common/config.Version=$(VERSION) -X google.golang.org/protobuf/reflect/protoregistry.conflictPolicy=warn" #-race

replenish: $(SRC_FILES)
	$(GC)  $(BUILD_NODE_PAR) -o replenish main.go

clean:
	rm -rf replenish