test: main.go
	go build .
	./mimic -s="./tests/.mimic" -t="./tests"