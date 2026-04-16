test.java: main.go
	go build -o ./bin/mimic main.go
	./bin/mimic \
		-s ./tests/java/.mimic \
		-t ./tests/java \
		-v domain=user \
		-v Domain=User
	./bin/mimic \
		-s ./tests/java/.mimic \
		-t ./tests/java \
		-v domain=product \
		-v Domain=Product
	./bin/mimic \
		-s ./tests/java/.mimic \
		-t ./tests/java \
		-v domain=cart \
		-v Domain=Cart

test.tsx: main.go
	go build -o ./bin/mimic main.go
	./bin/mimic \
		-s ./tests/tsx/.mimic \
		-t ./tests/tsx \
		-v domain=cart \
		-v Domain=Cart