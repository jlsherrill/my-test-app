FROM registry.redhat.io/ubi8/go-toolset:latest as builder

WORKDIR /go/src/app
COPY . .

USER 0

RUN go get -d ./... && \
    go build -o widgets  server.go

RUN cp /go/src/app/widgets /usr/bin/

FROM registry.redhat.io/ubi8/ubi-minimal:8.5

WORKDIR /

COPY --from=builder /go/src/app/widgets ./widgets

USER 1001

CMD ["/widgets"]
