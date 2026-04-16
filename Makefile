test.java: main.go
	go build -o ./bin/mimic main.go
	./bin/mimic \
		-s ./.test/java/.mimic \
		-t ./.test/java \
		-v domain=user \
		-v Domain=User
	./bin/mimic \
		-s ./.test/java/.mimic \
		-t ./.test/java \
		-v domain=product \
		-v Domain=Product
	./bin/mimic \
		-s ./.test/java/.mimic \
		-t ./.test/java \
		-v domain=cart \
		-v Domain=Cart

test.tsx: main.go
	go build -o ./bin/mimic main.go
	./bin/mimic \
		-s ./.test/tsx/.mimic \
		-t ./.test/tsx \
		-v domain=cart \
		-v Domain=Cart