FROM library/golang:1.10.2 AS build-env

WORKDIR /go/src/namerctl

ADD . /go/src/github.com/linkerd/namerctl

RUN go build -ldflags "-linkmode external -extldflags -static" -o /go/bin/namerctl /go/src/github.com/linkerd/namerctl/main.go

FROM scratch

COPY --from=build-env /go/bin/namerctl /namerctl

ENTRYPOINT ["/namerctl"]
