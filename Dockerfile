FROM library/golang:1.7.3

WORKDIR /go/src/namerctl

ADD . /go/src/github.com/buoyantio/namerctl

RUN go build -o /go/bin/namerctl /go/src/github.com/buoyantio/namerctl/main.go

ENTRYPOINT ["/go/bin/namerctl"]
