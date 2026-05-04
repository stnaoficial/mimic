VERSION := $(shell git tag --sort=-version:refname | head -n 1)
RELEASE_DIR := ./releases/$(VERSION)

build: main.go
	go build -o ./bin/mimic main.go

release: main.go
	mkdir -p "$(RELEASE_DIR)/linux-amd64"
	mkdir -p "$(RELEASE_DIR)/windows-amd64"
	mkdir -p "$(RELEASE_DIR)/darwin-amd64"

	GOOS=linux GOARCH=amd64 go build \
	-ldflags="-X main.Version=$(VERSION)" \
	-o "$(RELEASE_DIR)/linux-amd64/mimic" main.go

	GOOS=windows GOARCH=amd64 go build \
	-ldflags="-X main.Version=$(VERSION)" \
	-o "$(RELEASE_DIR)/windows-amd64/mimic.exe" main.go

	GOOS=darwin GOARCH=amd64 go build \
	-ldflags="-X main.Version=$(VERSION)" \
	-o "$(RELEASE_DIR)/darwin-amd64/mimic" main.go

	cp README.md "$(RELEASE_DIR)/linux-amd64"
	cp README.md "$(RELEASE_DIR)/windows-amd64"
	cp README.md "$(RELEASE_DIR)/darwin-amd64"

	cp MANUAL.md "$(RELEASE_DIR)/linux-amd64"
	cp MANUAL.md "$(RELEASE_DIR)/windows-amd64"
	cp MANUAL.md "$(RELEASE_DIR)/darwin-amd64"

	tar -czf "$(RELEASE_DIR)/mimic-linux-amd64.tar.gz" -C "$(RELEASE_DIR)" linux-amd64
	tar -czf "$(RELEASE_DIR)/mimic-windows-amd64.tar.gz" -C "$(RELEASE_DIR)" windows-amd64
	tar -czf "$(RELEASE_DIR)/mimic-darwin-amd64.tar.gz" -C "$(RELEASE_DIR)" darwin-amd64

	rm -rf "$(RELEASE_DIR)/linux-amd64"
	rm -rf "$(RELEASE_DIR)/windows-amd64"
	rm -rf "$(RELEASE_DIR)/darwin-amd64" 