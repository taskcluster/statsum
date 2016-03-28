
TAG := v1

install:
	go get github.com/tinylib/msgp
	go get github.com/pquerna/ffjson

generate:
	go generate ./...

build:
	go build && docker build -t jonasfj/statsum:latest .

push:
	docker tag jonasfj/statsum:latest jonasfj/statsum:${TAG}
	docker push jonasfj/statsum:${TAG}

test:
	go test -race -v ./...

bench:
	go test -bench . ./server -benchmem
