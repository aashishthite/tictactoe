

# All packages, except the vendor directory
all_go_packages := $(shell go list ./... | grep -v vendor)

test:
	go vet ${all_go_packages}
	go test ${all_go_packages} -v

build:
	go build -o ./bin/server ./cmd/server