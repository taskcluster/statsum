
install:
	go install github.com/tinylib/msgp/msgp

generate:
	go generate ./...

test:
	go test -race -v ./...

bench:
	go test -bench . ./server -benchmem
