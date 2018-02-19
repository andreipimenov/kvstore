FROM golang:1.10

RUN go get -u github.com/golang/dep/cmd/dep

COPY . /go/src/github.com/andreipimenov/kvstore

WORKDIR /go/src/github.com/andreipimenov/kvstore

RUN dep ensure

RUN go build -o /go/bin/server cmd/server/*.go

ENTRYPOINT ["/go/src/github.com/andreipimenov/kvstore/entrypoint.sh"]

CMD server --config=etc/server.conf.json
