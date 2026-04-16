test.java: main.go
	go build -o ./bin/mimic main.go
	./bin/mimic \
		-s ./tests/java/.mimic \
		-t ./tests/java \
		-v domain=user \
		-v class=UserEntity \
		-v interface=UserRepository

test.tsx: main.go
	go build -o ./bin/mimic main.go
	./bin/mimic \
		-s ./tests/tsx/.mimic \
		-t ./tests/tsx \
		-v domain=cart \
		-v component=CartComponent \
		-v interface=CartInterface