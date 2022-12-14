GOSRC = $(shell find . -name "*.go" ! -name "*fake*" ! -name "*test.go") version/VERSION

.PHONY: all
all: test build

.PHONY: build
build: pc

.PHONY: clean
clean:
	rm -f ./pc

.PHONY: install
install: pc
	go install .

.PHONY: lint
lint:
	golangci-lint run ./...

pc: $(GOSRC)
	go build -o pc  .

.PHONY: test
test: lint
	go run github.com/onsi/ginkgo/v2/ginkgo run --randomize-all --race ./...

.PHONY: watch-test
watch-test:
	go run github.com/onsi/ginkgo/v2/ginkgo watch ./...
