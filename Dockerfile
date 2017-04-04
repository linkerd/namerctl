FROM library/golang:1.7.3

WORKDIR /go/src/namerctl

ADD . /go/src/github.com/linkerd/namerctl

RUN go build -o /go/bin/namerctl /go/src/github.com/linkerd/namerctl/main.go

ENTRYPOINT ["/go/bin/namerctl"]
