test.java: main.go
	go build -o ./bin/mimic main.go
	./bin/mimic \
		-s ./example/java/.mimic \
		-t ./example/java \
		-v domain=user \
		-v Domain=User
	./bin/mimic \
		-s ./example/java/.mimic \
		-t ./example/java \
		-v domain=product \
		-v Domain=Product
	./bin/mimic \
		-s ./example/java/.mimic \
		-t ./example/java \
		-v domain=cart \
		-v Domain=Cart

test.tsx: main.go
	go build -o ./bin/mimic main.go
	./bin/mimic \
		-s ./example/tsx/.mimic \
		-t ./example/tsx \
		-v domain=cart \
		-v Domain=Cart