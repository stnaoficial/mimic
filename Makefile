build: main.go
	go build -o ./bin/mimic main.go

example.java: ./bin/mimic
	./bin/mimic \
		-s ./example/java/.mimic \
		-t ./example/java \
		-v domain="Foo Bar Baz"