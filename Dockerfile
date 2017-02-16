FROM        golang:1.7
MAINTAINER  Jonas Finnemann Jensen <jopsen@gmail.com>
ENV         PORT  443
EXPOSE      443

COPY        . /go/src/github.com/taskcluster/statsum
WORKDIR     /go/src/github.com/taskcluster/statsum
RUN         go get -v -d ./cmd/statsum
RUN         go install -v ./cmd/statsum

CMD         ["/go/bin/statsum"]
