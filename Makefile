build: main.go
	go build -o ./bin/mimic main.go

example.java: ./bin/mimic
	./bin/mimic \
		-v domain="Foo Bar Baz" \
		./example/java/.mimic \
		./example/java