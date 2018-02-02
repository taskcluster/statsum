
TAG := v9

install:
	go get github.com/pkg/errors
	go get github.com/tinylib/msgp
	go get github.com/pquerna/ffjson
	go get github.com/kardianos/govendor

generate:
	go generate ./...

build:
	govendor sync
	go build ./cmd/statsum && docker build -t jonasfj/statsum:latest .

push:
	docker tag jonasfj/statsum:latest jonasfj/statsum:${TAG}
	docker push jonasfj/statsum:${TAG}

test:
	go test -race -v ./...

bench:
	go test -bench . ./server -benchmem
