FROM golang:alpine AS builder

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

ADD . /go/src/byd50-ssi
WORKDIR /go/src/byd50-ssi
RUN go get byd50-ssi

COPY go.mod go.sum ./did-sep/main.go ./

RUN go mod download

RUN go build -o main .

WORKDIR /dist

RUN cp /go/src/byd50-ssi/main .

FROM scratch
COPY --from=builder /go/src/byd50-ssi/configs.yml .
COPY --from=builder /dist/main .

ENTRYPOINT ["/main"]