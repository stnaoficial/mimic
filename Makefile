# Testing project features
test: main.go
	go test -v ./tests/**

# Building project binaries
build: main.go
	go build -o ./bin/mimic main.go

# Releasing a new version
VERSION := $(shell git tag --sort=-version:refname | head -n 1)
RELEASE_DIR := ./releases/$(VERSION)

release: main.go
	mkdir -p "$(RELEASE_DIR)/linux-amd64"
	mkdir -p "$(RELEASE_DIR)/linux-arm64"
	mkdir -p "$(RELEASE_DIR)/darwin-amd64"
	mkdir -p "$(RELEASE_DIR)/darwin-arm64"

	GOOS=linux GOARCH=amd64 go build \
		-ldflags="-X main.Version=$(VERSION)" \
		-o "$(RELEASE_DIR)/linux-amd64/mimic" main.go

	GOOS=linux GOARCH=arm64 go build \
		-ldflags="-X main.Version=$(VERSION)" \
		-o "$(RELEASE_DIR)/linux-arm64/mimic" main.go

	GOOS=darwin GOARCH=amd64 go build \
		-ldflags="-X main.Version=$(VERSION)" \
		-o "$(RELEASE_DIR)/darwin-amd64/mimic" main.go

	GOOS=darwin GOARCH=arm64 go build \
		-ldflags="-X main.Version=$(VERSION)" \
		-o "$(RELEASE_DIR)/darwin-arm64/mimic" main.go

	cp MANUAL.md "$(RELEASE_DIR)/linux-amd64"
	cp MANUAL.md "$(RELEASE_DIR)/linux-arm64"
	cp MANUAL.md "$(RELEASE_DIR)/darwin-amd64"
	cp MANUAL.md "$(RELEASE_DIR)/darwin-arm64"

	tar -czf "$(RELEASE_DIR)/mimic-linux-amd64.tar.gz" \
		-C "$(RELEASE_DIR)" linux-amd64

	tar -czf "$(RELEASE_DIR)/mimic-linux-arm64.tar.gz" \
		-C "$(RELEASE_DIR)" linux-arm64

	tar -czf "$(RELEASE_DIR)/mimic-darwin-amd64.tar.gz" \
		-C "$(RELEASE_DIR)" darwin-amd64

	tar -czf "$(RELEASE_DIR)/mimic-darwin-arm64.tar.gz" \
		-C "$(RELEASE_DIR)" darwin-arm64

	rm -rf "$(RELEASE_DIR)/linux-amd64"
	rm -rf "$(RELEASE_DIR)/linux-arm64"
	rm -rf "$(RELEASE_DIR)/darwin-amd64"
	rm -rf "$(RELEASE_DIR)/darwin-arm64"