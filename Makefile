test.java: main.go
	go build -o ./bin/mimic main.go
	./bin/mimic \
		-s ./example/java/.mimic \
		-t ./example/java \
		-v Domain=FooBarBaz